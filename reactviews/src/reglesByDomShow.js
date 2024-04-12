import React, { useState, useEffect } from 'react';
import { useStore, useGetRecordId, useRecordContext, useDataProvider } from 'react-admin';
import { Show } from 'react-admin';

import { NavLink } from 'react-router-dom';
import Button from '@mui/material/Button';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';

import { RegleView } from './regles';
import { MyConfig } from './MyConfig';

const RegleTitle = () => {
    const record = useRecordContext();
    return <span>Règle {record ? `"${record.code}"` : ''}</span>;
};

const RegleByDomShow = () => {
    const recordId = useGetRecordId();
    const dataProvider = useDataProvider();
    //console.log(recordId);
    const [dom] = useStore('dom', {select: {id:0 }, evaluation: {}});
    const mask = dom.mask || { mask: {etu:false, nok:false} };
    const [state, setState] = useState({data: null, dom: 0});


    useEffect(() => {
        //console.log('-----------');
        dataProvider.getOne('regles/'+ recordId, {id: dom.select.id})
            .then(({ data }) => {
                //console.log('--dataProvider--');
                //console.log(data);
                setState({
                    data: data, etu: mask.etu, nok: mask.nok,
                });
            })
            .catch(error => {
                console.log('-error-', error);
                window.location.href = MyConfig.BASE_PATH + '#/login';
                return;
            });

    },[dom]);


    if (state.data === null) {
        //console.log(this.state);
        //console.log(this.dom);
        return <div>{state.data}</div>;
    } else {
        return (
            <Show title={<RegleTitle />} >

                <div style={{margin:'1em'}}>
                    <NavLink to={`/themes/${state.data.theme.id}/show`} style={{ textDecoration: 'none' }}>
                        <Button variant="outlined" color="primary" startIcon={<ArrowUpwardIcon />} >
                           Thème
                        </Button>
                    </NavLink> {state.data.theme.name}

                    {dom && dom.select.id !== '0' ?
                        <h2>Périmètre : {dom.select.name}</h2>
                        : null }
                    <RegleView record={state.data} />
                </div>
            </Show>
        );
    }

};

export default RegleByDomShow;
