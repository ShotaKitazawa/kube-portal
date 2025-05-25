/* eslint @next/next/no-img-element: 0 */

import React, { useState, useEffect, useContext, ReactNode } from 'react'
import { CircularProgress } from '@mui/material'

import { Links } from '../Links'
import { GlobalContext } from '../../contexts/global'
import { LinkInfo } from '../../drivers/ingress-info'

interface LinkGroupsProps {
  children?: ReactNode
}

export const LinkGroups: React.FC<LinkGroupsProps> = (props) => {
  const { linkInfo } = useContext(GlobalContext)
  const [linkInfoByTags, setLinkInfoByTags] =
    useState<Map<string, LinkInfo[]>>()

  const convertLinkInfo = (linkInfo: LinkInfo[]) => {
    let res: Map<string, LinkInfo[]> = new Map()
    for (let i = 0; i < linkInfo.length; i++) {
      for (let j = 0; j < linkInfo[i].tags.length; j++) {
        let tmp = res.get(linkInfo[i].tags[j]) ?? []
        tmp.push(linkInfo[i])
        res.set(linkInfo[i].tags[j], tmp.sort())
      }
    }
    setLinkInfoByTags(res)
  }

  useEffect(() => {
    if (linkInfo != null) {
      convertLinkInfo(linkInfo)
    }
  }, [linkInfo])

  return (
    <>
      {linkInfoByTags != null ? (
        Array.from(linkInfoByTags.keys())
          .sort()
          .map((tagName, index) => (
            <React.Fragment key={`tag-group-${index}`}>
              {tagName != '' ? (
                <React.Fragment key={`tag-header-${index}`}>
                  <div className="grid grid-cols-7">
                    <h2 className="col-start-2 text-4xl mt-20 font-bold">
                      {tagName}
                    </h2>
                  </div>
                  <div className="grid grid-cols-8">
                    <hr className="col-start-2 col-end-8 h-px my-4 bg-gray-200 border-0" />
                  </div>
                </React.Fragment>
              ) : (
                <React.Fragment key={`empty-tag-${index}`}></React.Fragment>
              )}
              <Links linksList={linkInfoByTags.get(tagName)} />
            </React.Fragment>
          ))
      ) : (
        <CircularProgress />
      )}
    </>
  )
}
