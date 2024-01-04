import { NextPage } from 'next'

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
