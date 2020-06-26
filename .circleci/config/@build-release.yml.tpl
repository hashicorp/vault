{{- $data := (datasource "package-list") -}}
{{- /*
  BUILD_ID is set by the staging workflow to produce an identifiable build.
*/ -}}
{{- $buildID := (env.Getenv "BUILD_ID" "standalone") -}}
{{- $workflowName := (env.Getenv "RELEASE_BUILD_WORKFLOW_NAME" "build-standalone") -}}
{{- $packages := $data.packages -}}
{{- $layers := $data.layers -}}
{{- $revision := (env.Getenv "PRODUCT_REVISION") -}}
{{- define "cache-key"}}{{template "cache-version"}}-{{.}}{{end -}}
{{- define "cache-version"}}cache002{{end -}}
{{- /*
  Any change to cache-version invalidates all build layer and package caches.
*/ -}}
# Current cache version: {{template "cache-version"}}

executors:
  releaser:
    docker:
      - image: circleci/buildpack-deps
    environment:
      PRODUCT_REVISION: "{{if $revision}}{{$revision}}{{end}}"
      AUTO_INSTALL_TOOLS: 'YES'
    shell: /usr/bin/env bash -euo pipefail -c

workflows:
  {{$workflowName}}:
    jobs:
      - cache-builder-images:
          filters:
            branches:
              only:
                - /build-.*/
                - /ci.*/
      {{- range $packages}}
      - {{.meta.BUILD_JOB_NAME}}: { requires: [ cache-builder-images ] }
      {{- end }}
      - bundle-releases:
          requires:
            {{- range $packages}}
            - {{.meta.BUILD_JOB_NAME}}{{end}}
jobs:
  cache-builder-images:
    executor: releaser
    steps:
      - setup_remote_docker
      - checkout
      - write-build-layer-cache-keys

      # Load best available cached image.

      {{- $targetLayerType := "build-static-assets" }}
      {{- $targetLayer := .}}
      {{- range $layers}}
      {{- if eq .type $targetLayerType }}
      {{- $targetLayer = .}}
      - restore_cache:
          keys:
            {{- range .meta.circleci.CACHE_KEY_PREFIX_LIST}}
            - {{template "cache-key" .}}
            {{- end}}
      - run:
          name: Finish early if loaded exact match from cache.
          command: |
            if [ -f {{.archivefile}} ]; then
              echo "Exact match found in cache, skipping build."
              circleci-agent step halt
            else
              echo "No exact match found, proceeding with build."
            fi
      - run: LAYER_SPEC_ID={{.name}} make -C packages*.lock load-builder-cache
      {{- end}}{{end}}

      # No exact match was found, so build each layer up to target type.

      {{- $finished := false }}
      {{- range $layers}}
      {{- if not $finished }}
      {{- $finished = eq .type $targetLayerType}}
      - run: make -f packages*.lock/layer.mk {{.name}}-image
      {{- end}}
      {{- end}}

      # Save the target layer archive.
      
      - run: make -f packages*.lock/layer.mk {{$targetLayer.name}}-save
      
      # Store the target layer archive as all the relevant cache names.

      {{- $lastArchive := $targetLayer.archivefile}}
      {{- range $i, $l := $targetLayer.meta.builtin.LAYER_LIST}}
      {{- $currentArchive := $l.archive}}
      {{- if ne $currentArchive $lastArchive }}
      - run: mv {{$lastArchive}} {{$currentArchive}}
      {{- end}}
      - save_cache:
          key: {{template "cache-key" (index $targetLayer.meta.circleci.CACHE_KEY_PREFIX_LIST $i)}}
          paths:
            - {{$currentArchive}}
      {{- $lastArchive = $currentArchive }}
      {{- end}}

{{- range $packages}}
  {{.meta.BUILD_JOB_NAME}}:
    executor: releaser
    environment:
      - PACKAGE_SPEC_ID: {{.packagespecid}}
    steps:
      - setup_remote_docker
      - checkout

      # Restore the package cache first, we might not need to rebuild.
      - write-package-cache-key
      - restore_cache:
          key: '{{template "cache-key" .meta.circleci.PACKAGE_CACHE_KEY}}'
      - run:
          name: Check the cache status.
          command: |
            if ! { PKG=$(find .buildcache/packages/store -maxdepth 1 -mindepth 1 -name '*.zip' 2> /dev/null) && [ -n "$PKG" ]; }; then
              echo "No package found, continuing with build."
              exit 0
            fi
            echo "Package already cached, skipping build."
            circleci-agent step halt

      # We need to rebuild, so load the builder cache.
      - write-build-layer-cache-keys
      - restore_cache:
          keys:
          {{- range .meta.circleci.BUILDER_CACHE_KEY_PREFIX_LIST}}
          - {{template "cache-key" .}}
          {{- end}}
      - run: make -C packages*.lock load-builder-cache
      - run: make -C packages*.lock package
      - run: ls -lahR .buildcache/packages
      # Save package cache.
      - save_cache:
          key: '{{template "cache-key" .meta.circleci.PACKAGE_CACHE_KEY}}'
          paths:
            - .buildcache/packages/store
      # Save builder image cache if necessary.
      # The range should only iterate over a single layer.
      {{- $pkg := . -}}
      {{- range $idx, $layerInfo := .meta.builtin.BUILD_LAYERS }}
      {{- if eq $layerInfo.type "warm-go-build-vendor-cache" }}
      {{- with $layerInfo }}
      {{- $circleCICacheKey := (index $pkg.meta.circleci.BUILDER_CACHE_KEY_PREFIX_LIST $idx) }}
      - run:
          name: Check builder cache status
          command: |
            if [ -f {{.archive}} ]; then
              echo "Builder image already cached, skipping cache step."
              circleci-agent step halt
            fi
      - run: make -f packages*.lock/layer.mk {{.name}}-save
      - save_cache:
          key: '{{template "cache-key" $circleCICacheKey}}'
          paths:
            - {{.archive}}
      {{- end}}
      {{- end}}
      {{- end}}
{{end}}

  bundle-releases:
    executor: releaser
    steps:
      - checkout
      - write-all-package-cache-keys
      {{- range $packages}}
      - load-{{.meta.BUILD_JOB_NAME}}
      - run:
          environment:
            PACKAGE_SPEC_ID: {{.packagespecid}}
          name: Write package metadata for {{.meta.BUILD_JOB_NAME}}
          command: |
            make package-meta
      {{- end}}
      - run:
          name: Write package aliases
          command:
            make aliases
      - run:
          name: List Build Cache
          command: ls -lahR .buildcache

      # Surface the package store directory as an artifact.
      # This makes each zipped package separately downloadable.
      - store_artifacts:
          path: .buildcache/packages
          destination: packages-{{$buildID}}

      # Surface a tarball of the whole package store as an artifact.
      - run: tar -czf packages-{{$buildID}}.tar.gz .buildcache/packages
      - store_artifacts:
          path: packages-{{$buildID}}.tar.gz
          destination: packages-{{$buildID}}.tar.gz

      # Surface a tarball of just the metadata files.
      - run: tar -czf meta-{{$buildID}}.tar.gz .buildcache/packages/store/*.json
      - store_artifacts:
          path: meta-{{$buildID}}.tar.gz
          destination: meta-{{$buildID}}.tar.gz

commands:
  {{- range $packages }}
  load-{{.meta.BUILD_JOB_NAME}}:
    steps:
      - restore_cache:
          key: '{{template "cache-key" .meta.circleci.PACKAGE_CACHE_KEY}}'
  {{end}}

  write-build-layer-cache-keys:
    steps:
      - run:
          name: Write builder layer cache keys
          command: make -C packages*.lock write-builder-cache-keys

  write-package-cache-key:
    steps:
      - run:
          name: Write package cache key
          command: make -C packages*.lock write-package-cache-key

  write-all-package-cache-keys:
    steps:
      - run:
          name: Write package cache key
          command: make -C packages*.lock write-all-package-cache-keys
