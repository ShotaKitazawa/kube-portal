import { NextPage } from 'next'

import { Layout } from '../components/Layout'
import { Links } from '../components/Links'

const IndexPage: NextPage = (props) => {
  return (
    <Layout>
      <Links />
    </Layout>
  )
}

export default IndexPage
