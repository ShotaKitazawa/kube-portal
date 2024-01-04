/* eslint @next/next/no-img-element: 0 */

import { useState, useEffect, useContext } from 'react'
import { CircularProgress } from '@mui/material'

import { Links } from '../Links'
import { GlobalContext } from '../../contexts/global'
import { LinkInfo } from '../../drivers/ingress-info'

type Props = {}

export const LinkGroups: React.FC<Props> = (props) => {
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
        Array.from(linkInfoByTags.keys()).map((tagName) => (
          <>
            {tagName != '' ? (
              <>
                <div className="grid grid-cols-7">
                  <h2 className="col-start-2 text-4xl mt-20 font-bold">
                    {tagName}
                  </h2>
                </div>
                <div className="grid grid-cols-8">
                  <hr className="col-start-2 col-end-8 h-px my-4 bg-gray-200 border-0" />
                </div>
              </>
            ) : (
              <></>
            )}
            <Links linksList={linkInfoByTags.get(tagName)} />
          </>
        ))
      ) : (
        <CircularProgress />
      )}
    </>
  )
}
