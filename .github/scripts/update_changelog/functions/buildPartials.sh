#!/bin/bash

function buildPartials {

  if [[ -z "${1}" ]] ; then return ; fi

  local TAG_KEY="<CURR_TAG>"
  local VERSION_N="${1}"  
  local VERSION_NM1=$( awk '{print $1-$2}' <<< "${VERSION_N} ${MAJOR_VERSION_DELTA}" )
  local VERSION_NM2=$( awk '{print $1-$2}' <<< "${VERSION_NM1} ${MAJOR_VERSION_DELTA}" )

  buildMarkdown ${VERSION_N}
  buildMarkdown ${VERSION_NM1}
  buildMarkdown ${VERSION_NM2}

  mv "${OUTPUT_DIR}/changelog_${VERSION_N}.mdx"   "${OUTPUT_DIR}/version-n.mdx"
  mv "${OUTPUT_DIR}/changelog_${VERSION_NM1}.mdx" "${OUTPUT_DIR}/version-nm1.mdx"
  mv "${OUTPUT_DIR}/changelog_${VERSION_NM2}.mdx" "${OUTPUT_DIR}/version-nm2.mdx"

  cat "${OUTPUT_DIR}/version-n.mdx" | sed s/"${TAG_KEY}"/" (latest)"/ > "${OUTPUT_DIR}/version-n.mdx"
  cat "${OUTPUT_DIR}/version-nm1.mdx" | sed s/"${TAG_KEY}"/""/ > "${OUTPUT_DIR}/version-nm1.mdx"
  cat "${OUTPUT_DIR}/version-nm2.mdx" | sed s/"${TAG_KEY}"/""/ > "${OUTPUT_DIR}/version-nm2.mdx"
}

function buildMarkdown {

   if [[ -z "${1}" ]] ; then return ; fi

  local VERSION_N="${1}"  
  local VERSION_NM1=$( awk '{print $1-$2}' <<< "${VERSION_N} ${MAJOR_VERSION_DELTA}" )

  local H2_REGEX="^## .*"
  local H3_REGEX="^### .*"
  local VERSION_REGEX="${H2_REGEX}${VERSION_N}.*"
  local STOP_REGEX="${H2_REGEX}${VERSION_NM1}.*"
  local DATE_REGEX="${H3_REGEX}[A-Za-z].*"
  local CAPS_HEADER_REGEX="^[A-Z ]+: *"

  local STOP=0
  local NEW_ENTRY=0
  local START_COPY=0
  local WRITE_LINE=0
  local LIST_ITEM=0

  local COPY_LINE=""
  local TOC_LINE=""
  local TRIMMED_DATE=""
  local ANCHOR="" 
  local VERSION_KEY="<VERSION>"
  local TOC_KEY="<TOC>"
  local DETAILS_KEY="<DETAILS>"
  local V_PARTIAL="${OUTPUT_DIR}/changelog_${VERSION_N}.mdx"
  local V_CHANGELOG="${OUTPUT_DIR}/body_${VERSION_N}.mdx"
  local V_TOC="${OUTPUT_DIR}/toc_${VERSION_N}.mdx"

  # Initialize the scratch files
  cat "${MD_TAB_BODY}" | sed s/${VERSION_KEY}/${VERSION_N}/ > "${V_CHANGELOG}" # Scratch file for changelog entries
  cat "${MD_TOC}" | sed s/${VERSION_KEY}/${VERSION_N}/ > "${V_TOC}"            # Scratch file for tab TOC

  # Copy lines from the changelog until we find an entry for the previous version
  while read -r CURR_LINE
  do

    # The current line is a version header. Flip the start copy and new
    # entry values and set copy line to the current line
    if [[ "${CURR_LINE}" =~ ${VERSION_REGEX} ]]
    then
      START_COPY=1
      NEW_ENTRY=2
      COPY_LINE="${CURR_LINE}"
      ANCHOR=${CURR_LINE//#/}
      ANCHOR=${ANCHOR//[[:space:]]/}
      TOC_LINE=""

    # The current line is a date entry. Append it to the version header
    elif [[ "${CURR_LINE}" =~ ${DATE_REGEX} ]]
    then
      TRIMMED_DATE=$( echo ${CURR_LINE//#/} | sed 's/^[ \t]*//;s/[ \t]*$//')
      TOC_LINE="- [${ANCHOR} — ${TRIMMED_DATE}](#${ANCHOR})"
      COPY_LINE="${COPY_LINE} — ${TRIMMED_DATE} ((#${ANCHOR}))"
      NEW_ENTRY=0

    # The current line is a version header for the previous version. Stop reading.
    elif [[ "${CURR_LINE}" =~ ${STOP_REGEX} ]]
    then
      STOP=1

    # The current line is a section header (e.g., "SECURITY, IMPROVEMENTS")
    # Remove the colon, convert the line to sentence case, and make it an h3
    elif [[ "${CURR_LINE}" =~ ${CAPS_HEADER_REGEX} ]]
    then
      COPY_LINE="${CURR_LINE//:/}"
      COPY_LINE="${COPY_LINE,,}"
      COPY_LINE="### ${COPY_LINE^}"
      TOC_LINE=""
    
    # The current line is a change entry, replace the * list item marker with -
    elif [[ "${CURR_LINE:0:1}" == "*" ]]
    then
      LIST_ITEM=3
      COPY_LINE="-${CURR_LINE:1}"
      TOC_LINE=""
    
    # The current line is (potentially) the continuation of a change entry
    elif [[ ! -z "${CURR_LINE}" ]]
    then
      COPY_LINE=$( echo ${CURR_LINE} | sed 's/^[ \t]*//;s/[ \t]*$//')
      COPY_LINE="  ${COPY_LINE}"
      LIST_ITEM=0
      TOC_LINE=""

    # The current line is (probably) empty
    else
      COPY_LINE=${CURR_LINE}
      TOC_LINE=""
    fi

    WRITE_LINE=$((START_COPY + NEW_ENTRY))

    if [[ ${STOP} -eq 1 ]] ; then break ; fi
    if [[ ${WRITE_LINE} -eq 1 ]] ; then

      echo "${COPY_LINE}" >> "${V_CHANGELOG}"
      
      if [[ -n "${TOC_LINE}" ]] ; then
        echo "${TOC_LINE}" >> "${V_TOC}"
      fi
    fi

  done < ${LOCAL_CHANGELOG}


  local V_PARTIAL="${OUTPUT_DIR}/changelog_${VERSION_N}.mdx"
  local V_CHANGELOG="${OUTPUT_DIR}/body_${VERSION_N}.mdx"
  local V_TOC="${OUTPUT_DIR}/toc_${VERSION_N}.mdx"

  # Write the final partial file
  cat "${MD_TAB_OPEN}" | sed s/"${VERSION_KEY}"/"${VERSION_N}.x"/ > "${V_PARTIAL}"
  cat "${V_TOC}"        >> "${V_PARTIAL}"
  echo ""               >> "${V_PARTIAL}"
  cat "${V_CHANGELOG}"  >> "${V_PARTIAL}"
  echo ""               >> "${V_PARTIAL}"
  cat "${MD_TAB_CLOSE}" >> "${V_PARTIAL}"

  ## Cleanup
  rm "${V_TOC}"
  rm "${V_CHANGELOG}"
}