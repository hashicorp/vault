"""This is a modified version of snapcraft's go plugin.

It uses go 1.7 from ppa:gophers/archive, because it is not yet published
in the xenial archives.

The go plugin can be used for go projects using `go get`.

This plugin uses the common plugin keywords, for more information check the
'plugins' topic.

This plugin uses the common plugin keywords as well as those for "sources".
For more information check the 'plugins' topic for the former and the
'sources' topic for the latter.

Additionally, this plugin uses the following plugin-specific keywords:

    - go-importpath:
      (string)
      This entry tells the checked out `source` to live within a certain path
      within `GOPATH`.
    - go-buildtags:
      (list of strings)
      Tags to use during the go build. Default is not to use any build tags.

"""

import logging
import os
import shutil

import snapcraft
from snapcraft import common


logger = logging.getLogger(__name__)


class GoPlugin(snapcraft.BasePlugin):

    @classmethod
    def schema(cls):
        schema = super().schema()
        schema['properties']['go-importpath'] = {
            'type': 'string',
            'default': ''
        }
        schema['properties']['go-buildtags'] = {
            'type': 'array',
            'minitems': 1,
            'uniqueItems': True,
            'items': {
                'type': 'string',
            },
            'default': []
        }

        schema['required'].append('go-importpath')

        # Inform Snapcraft of the properties associated with pulling. If these
        # change in the YAML Snapcraft will consider the pull step dirty.
        schema['pull-properties'].append('go-importpath')

        # Inform Snapcraft of the properties associated with building. If these
        # change in the YAML Snapcraft will consider the build step dirty.
        schema['build-properties'].extend(['source', 'go-importpath'])

        return schema

    def __init__(self, name, options, project):
        super().__init__(name, options, project)
        self.build_packages.append('golang-1.7-go')
        self._gopath = os.path.join(self.partdir, 'go')
        self._gopath_src = os.path.join(self._gopath, 'src')
        self._gopath_bin = os.path.join(self._gopath, 'bin')
        self._gopath_pkg = os.path.join(self._gopath, 'pkg')

    @property
    def go_bin(self):
        return '/usr/lib/go-1.7/bin/go'

    def pull(self):
        # use -d to only download (build will happen later)
        # use -t to also get the test-deps
        # since we are not using -u the sources will stick to the
        # original checkout.
        super().pull()
        os.makedirs(self._gopath_src, exist_ok=True)

        go_package = self.options.go_importpath
        go_package_path = os.path.join(self._gopath_src, go_package)
        if os.path.islink(go_package_path):
            os.unlink(go_package_path)
        os.makedirs(os.path.dirname(go_package_path), exist_ok=True)
        os.symlink(self.sourcedir, go_package_path)
        self._run(
            [self.go_bin, 'get', '-t', '-d', './{}/...'.format(go_package)])

    def clean_pull(self):
        super().clean_pull()

        # Remove the gopath (if present)
        if os.path.exists(self._gopath):
            shutil.rmtree(self._gopath)

    def build(self):
        super().build()

        tags = []
        if self.options.go_buildtags:
            tags = ['-tags={}'.format(','.join(self.options.go_buildtags))]
        self._run([self.go_bin, 'install'] + tags +
                  ['./{}/...'.format(self.options.go_importpath)])

        install_bin_path = os.path.join(self.installdir, 'bin')
        os.makedirs(install_bin_path, exist_ok=True)
        os.makedirs(self._gopath_bin, exist_ok=True)
        for binary in os.listdir(os.path.join(self._gopath_bin)):
            binary_path = os.path.join(self._gopath_bin, binary)
            shutil.copy2(binary_path, install_bin_path)

    def clean_build(self):
        super().clean_build()

        if os.path.isdir(self._gopath_bin):
            shutil.rmtree(self._gopath_bin)

        if os.path.isdir(self._gopath_pkg):
            shutil.rmtree(self._gopath_pkg)

    def _run(self, cmd, **kwargs):
        env = self._build_environment()
        return self.run(cmd, cwd=self._gopath_src, env=env, **kwargs)

    def _build_environment(self):
        env = os.environ.copy()
        env['GOPATH'] = self._gopath

        include_paths = []
        for root in [self.installdir, self.project.stage_dir]:
            include_paths.extend(
                common.get_library_paths(root, self.project.arch_triplet))

        flags = common.combine_paths(include_paths, '-L', ' ')
        env['CGO_LDFLAGS'] = '{} {} {}'.format(
            env.get('CGO_LDFLAGS', ''), flags, env.get('LDFLAGS', ''))

        return env
