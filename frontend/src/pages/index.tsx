import type { NextPage } from 'next'
import React from 'react'

import { Layout } from '../components/Layout'
import { LinkGroups } from '../components/LinkGroups'

const IndexPage: NextPage = () => {
  return (
    <Layout>
      <LinkGroups />
    </Layout>
  )
}

export default IndexPage
