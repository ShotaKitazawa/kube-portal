import {
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  CircularProgress,
} from '@material-ui/core';
import LaunchIcon from '@material-ui/icons/Launch';

import IngressInfo from '../../drivers/ingress-info/ingress-info'
import { LinksPort } from '../../entities/ingress-info'

export const Links: React.FC = ({
  children,
}) => {

  // get links from API
  const links: LinksPort = new IngressInfo()
  const linksList = links.List()
  const item = (name: string, icon_url: string) => (
    <ListItem>
      <ListItemAvatar>
        <Avatar>
          {
            icon_url != "" ?
              <img src={icon_url} />
              :
              <LaunchIcon />
          }
        </Avatar>
      </ListItemAvatar>
    </ListItem>
  )

  return (
    <>
      {
        linksList != null ?
          linksList.map(({ name, url, icon_url, is_private }) => (
            <>
              {
                <a href={url} className="hover:no-underline">
                  {item(name, icon_url)}
                </a>
              }
            </>
          ))
          :
          <CircularProgress />
      }
    </>
  )
}
