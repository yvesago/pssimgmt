import React from 'react';
import { MyConfig } from './MyConfig';

const authLoginPage  = () => {

    return (
        <div>
            <h1>{MyConfig.NAME}</h1>
            <p>Login with <a href={MyConfig.AUTH_URL}>CAS</a></p>
        </div>
    );
};

export default authLoginPage;
