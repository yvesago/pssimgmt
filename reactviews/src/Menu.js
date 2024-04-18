import React from 'react';
import { Logout, MenuItemLink, usePermissions, useResourceDefinitions } from 'react-admin';

import Divider from '@mui/material/Divider';
import ViewListSharpIcon from '@mui/icons-material/ViewListSharp';
import ListIcon from '@mui/icons-material/List';
import UserIcon from '@mui/icons-material/People';
import GroupWorkRoundedIcon from '@mui/icons-material/GroupWorkRounded';
import HistoryRoundedIcon from '@mui/icons-material/HistoryRounded';
import LocalLibrarySharpIcon from '@mui/icons-material/LocalLibrarySharp';
import AccountTreeSharpIcon from '@mui/icons-material/AccountTreeSharp';
import LibraryBooksSharpIcon from '@mui/icons-material/LibraryBooksSharp';
import FactCheckIcon from '@mui/icons-material/FactCheck';

const allMenu = ['home', 'todos', 'regles', 'versions'];
const icons = { 
    'users': <UserIcon />, 
    'versions': <HistoryRoundedIcon />,  
    'domaines': <GroupWorkRoundedIcon />,
    'docs': <LocalLibrarySharpIcon />,
    'themes': <AccountTreeSharpIcon />,
    'todos': <FactCheckIcon />,
    'documents': <LibraryBooksSharpIcon />,
};

const Menu = ({ onMenuClick }) => {
    const { permissions } = usePermissions();
    const resourcesDefinitions = useResourceDefinitions();
    const resources = Object.keys(resourcesDefinitions).map(name => resourcesDefinitions[name]);

    return (
        <div>
            <MenuItemLink key='home' primaryText='PSSI' leftIcon={<ViewListSharpIcon />} to='/' />
            { resources.map( (resource) =>  { 
                if (allMenu.includes(resource.name)) {
                    return (<MenuItemLink
                        key={resource.name}
                        primaryText= {resource.options && (resource.options.label || resource.name)}
                        to={`/${resource.name}`}
                        leftIcon={icons[resource.name] ? icons[resource.name]: <ListIcon /> }
                        onClick={onMenuClick}
                    />
                    );
                }
                else  { return null;}
            }).sort()
            }
            <Logout />
            <Divider />
            { (permissions === 'admin' ) ? 'Admin': '' }
            { resources.map( (resource) =>  { 
                if (allMenu.includes(resource.name) == false && permissions === 'admin' ) {
                    return (<MenuItemLink
                        key={resource.name}
                        primaryText= {resource.options && (resource.options.label || resource.name)}
                        to={`/${resource.name}`}
                        leftIcon={icons[resource.name] ? icons[resource.name]: <ListIcon /> }
                        onClick={onMenuClick}
                    />
                    );
                } else  { return null;}
            })
            }
        </div>
    );
};

export default Menu;
