import * as React from 'react';

import jsonServerProvider from 'ra-data-json-server';
import { fetchUtils, Admin, Resource, CustomRoutes, defaultTheme } from 'react-admin';
const lightTheme = defaultTheme;
const darkTheme = { ...defaultTheme, palette: { mode: 'dark' } };

import { Route } from 'react-router-dom';
import { ThemeList, ThemeEdit, ThemeCreate } from './themes';
import ThemeShow from './themesShow';
import { RegleList, RegleEdit, RegleCreate, RegleCodeShow } from './regles';
import RegleShow from './reglesShow';
import { UserList, UserEdit, UserShow } from './users';
import { DomaineList, DomaineEdit, DomaineCreate } from './domaines';
import { DocList, DocEdit, DocCreate } from './doc';
import { IsoList, IsoEdit, IsoCreate } from './iso';
import { VersionList, VersionEdit, VersionCreate } from './version';
import { DocumentList, DocumentEdit, DocumentCreate, DocumentShow } from './document';
import { TodoList } from './todo';
import { About } from './about';

import HomeShow from './home';

import UserIcon from '@mui/icons-material/People';

import authLoginPage from './authLoginPage';
import authClient from './authClient';


import Menu from './Menu';

import { MyConfig } from './MyConfig';

const httpClient = (url, options = {}) => {
    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }
    //options.headers.set('Access-Control-Allow-Headers', 'Authorization, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Origin, Accept, Content-Type, Location, Vary');
    // add your own headers here
    //options.headers.set('X-MyToken', MyConfig.API_KEY );
    const token = localStorage.getItem('ttoken');
    //console.log('httpClient == ' + token);
    if (token !== null) {
        options.headers.set('Authorization', `Bearer ${token}`);
    }
    options.headers.set('Access-Control-Allow-Origin', '*');
    options.headers.set('Access-Control-Allow-Methods', 'OPTIONS, GET, PUT, POST, DELETE');
    options.headers.set('Content-Type', 'application/json');
    options.credentials = 'include';
    return fetchUtils.fetchJson(url.replace('v1//domainestree','v1/domainestree'), options); //regex to request domainestree with dataprovider
};

const dataProvider = jsonServerProvider( MyConfig.API_URL , httpClient);

const App = () => (
    <Admin
        title={MyConfig.NAME}
        theme={lightTheme}
        darkTheme={darkTheme}
        menu={Menu}
        loginPage={authLoginPage}
        authProvider={authClient(MyConfig.AUTH_URL)}
        dashboard={HomeShow}
        dataProvider={dataProvider} disableTelemetry>
        {permissions => (
            <>
                <CustomRoutes noLayout>
                    <Route path="/themestree" element={<ThemeList />} />,
                </CustomRoutes>,
                <Resource name="users" options={{ label: 'Users'}} list={UserList} edit={UserEdit} show={UserShow} icon={UserIcon} />,
                <Resource options={{ label:'Périmètres'}}  name="domaines" list={DomaineList} edit={DomaineEdit} create={DomaineCreate} />,
                <Resource options={{ label:'Thèmes'}} name="themes" list={ThemeList} edit={ThemeEdit} show={ThemeShow} create={ThemeCreate} />,
                <Resource options={{ label:'Règle'}} name="regle" edit={RegleCodeShow} />,
                <Resource 
                    name="regles" 
                    options={{ label:'Règles'}} 
                    list={RegleList} 
                    edit={permissions === 'admin' ? RegleEdit: null} 
                    show={RegleShow} 
                    create={permissions === 'admin' ? RegleCreate: null} />,
                {permissions === 'admin' || permissions === 'cssi' ? <Resource options={{ label:'TODO'}} name="todos" list={TodoList} /> : null },
                <CustomRoutes>
                    <Route path="/about" element={<About />} />
                </CustomRoutes>
                <Resource 
                    name="versions" 
                    options={{ label:'Versions'}} 
                    list={VersionList} 
                    edit={permissions === 'admin' ? VersionEdit: null} 
                    create={permissions === 'admin' ? VersionCreate: null} />,
                <Resource options={{ label:'Refs. docs'}} name="docs" list={DocList} edit={DocEdit} create={DocCreate} />,
                <Resource options={{ label:'ISOs 27002'}} name="isos" list={IsoList} edit={IsoEdit} create={IsoCreate} />,
                <Resource options={{ label:'Documents'}} name="documents" list={DocumentList} edit={DocumentEdit} create={DocumentCreate} show={DocumentShow} />,
            </>
        )}
    </Admin>
);

export default App;
