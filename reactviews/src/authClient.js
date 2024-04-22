import decodeJwt from 'jwt-decode';
import request from 'superagent';
import { MyConfig } from './MyConfig';
import { PreviousLocationStorageKey } from 'react-admin';

const authClient = (authUrl) => ({
    login: async () => {return request.get(authUrl);},
    logout: async () => {
        localStorage.removeItem('ttoken');
        localStorage.removeItem(PreviousLocationStorageKey);
        return Promise.resolve('/login');
    },
    //checkAuth: () => Promise.resolve(),
    checkAuth: async () => {
        console.log('**checkAuth**',window.location.hash);
        const url = new URL(window.location.href);
        const match = url.hash.match(/\/login$/);
        const match2 = url.hash.match(/^#\/(\w+\/?)(\d+\/?)(\w+?)$/);
        if ( url.hash !== '' && match === null && match2 !== null ) {
            localStorage.setItem(PreviousLocationStorageKey, url.href);
        }

        if (localStorage.getItem('ttoken') !== null) {
            //console.log('**checkAuth**: auth', match2);
            return Promise.resolve();
        } else {
            //console.log('**checkAuth**: Unauth');
            return Promise.reject();
        }
    },
    checkError: async (error) => {
        const status = error.status;
        if (status === 401 || status === 403) {
            //console.log(error);
            localStorage.removeItem(PreviousLocationStorageKey);
            return Promise.reject({ message: 'Utilisateur non autorisé. En attente de validation.' });
        }
        return Promise.resolve();
    },
    getIdentity: async () => {
        if ( localStorage.getItem('ttoken') !== null ) {
            const decodedToken = decodeJwt(localStorage.getItem('ttoken'));
            //console.log('== getIdentity() : ' + JSON.stringify({ id: decodedToken.IDuser, fullName: decodedToken.Name, avatar: '' }) );
            return { id: decodedToken.IDuser, fullName: decodedToken.Name, avatar: '' };
        }
        return { id: '', fullName: '', avatar: '' };
    },
    //getPermissions: params => Promise.resolve(),
    getPermissions: async () => {
        if ( localStorage.getItem('ttoken') !== null ) {
            const decodedToken = decodeJwt(localStorage.getItem('ttoken'));
            const role = decodedToken.Role;
            //console.log('== Role: ' + role );
            return role ? Promise.resolve(role) : Promise.reject();
        }
        return Promise.reject();
    },
    handleCallback: async () => {
        console.log('**handleCallback**');
        var match = window.location.href.match(/\?(.*)$/);
        //console.log(match[1]);
        const token = match[1];
        localStorage.setItem('ttoken', token);
        //console.log(token);
        //window.location.href = MyConfig.BASE_PATH;
        window.location.href = localStorage.getItem(PreviousLocationStorageKey) !== null ? localStorage.getItem(PreviousLocationStorageKey) : MyConfig.BASE_PATH;
        return Promise.resolve();
    },
});

export default authClient;
