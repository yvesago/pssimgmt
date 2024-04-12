import simpleRestProvider from 'ra-data-simple-rest';
//import jsonServerProvider from 'ra-data-json-server';
//import { cacheDataProviderProxy } from 'react-admin'; 

const dataProvider = simpleRestProvider('http://127.0.0.:8000/');
//const dataProvider = jsonServerProvider('http://127.0.0.:8000/');


const cacheDataProviderProxy = (dataProvider, duration =  1) =>
    new Proxy(dataProvider, {
        get: (target, name) => (resource, params) => {
            if (name === 'getOne' || name === 'getMany' || name === 'getList') {
                return dataProvider[name](resource, params).then(response => {
                    const validUntil = new Date();
                    validUntil.setTime(validUntil.getTime() + duration);
                    response.validUntil = validUntil;
                    return response;
                });
            }
            return dataProvider[name](resource, params);
        },
    });


export default cacheDataProviderProxy(dataProvider);
