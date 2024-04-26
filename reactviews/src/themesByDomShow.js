import React, { useState, useEffect } from 'react';
import { useStore, usePermissions, useGetRecordId, useRedirect } from 'react-admin';
import { Title, Form, EditButton, useDataProvider, useTheme } from 'react-admin';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';

import FormGroup from '@mui/material/FormGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import Switch from '@mui/material/Switch';

import Divider from '@mui/material/Divider';

import { RegleView } from './regles';

import ReactMarkdown from 'react-markdown';


export const RegleMaskForm = () => {
    const [dom,setDom] = useStore('dom', { select: {id:0 }, evaluation: {}, mask: {etu:false, nok:false} });
    const mask = dom.mask || { mask: {etu:false, nok:false} }; 
    const [state, setState] = useState({ etu: mask.etu, nok: mask.nok });


    const handleChange = (event) => {
        setState({ ...state, [event.target.name]: event.target.checked });
        mask[event.target.name] = event.target.checked;
        setDom({select:dom.select, evaluation: dom.evaluation, mask: mask});
    };

    return (
        <>
            <Divider />
            <Form defaultValues={dom.mask}>
                <FormLabel component="legend">Visualisation des règles:</FormLabel>
                <FormGroup row>
                    <FormControlLabel
                        control={<Switch checked={state.etu} onChange={handleChange} name="etu" />}
                        label="En étude"
                    />
                    <FormControlLabel
                        control={<Switch checked={state.nok} onChange={handleChange} name="nok" />}
                        label="Non ok"
                    />
                </FormGroup>
            </Form>
            <Divider />
        </>
    );

};

const renderRegles = (r, mask) => {
    //console.log('render regles');
    //console.log(r);
    if (r.record !== null && (r.status === 'ok' || r.status === '' || (r.status === 'etu' && mask.etu) ||  (r.status === 'nok' && mask.nok ) )) {
        return(
            <li key={r.id}>
                <RegleView record={r} viewonly={false} /> 
            </li>
        );
    }
    return null;
};


const box = (axe, value, index, theme) => (
    <Box width={value} bgcolor={theme === 'dark'?'grey':'grey.300'} p={0.7} my={0.5} key={index}>
        {axe}: {value}
    </Box>
);


const RenderEval = (e) => (
    <Box key={e.record.id}>
  Conformité des règles sur les axes :
        {  e.record.axes.map( (axe, index) => {
            if ( axe !== '' ) {
                const value = 10 * e.record.conforme[index];
                if ( value > 10 ) {
                    return box(axe, value + '%', index, e.theme);
                }
                else {
                    return <Box key={index}>{axe}: {value}%</Box>;
                }
            }
            return null;
        })
        } 
    </Box>
);

const ThemesButton = () => {
    const redirect = useRedirect();
    const handleClick = () => {
        redirect('/');
    };
    return <Button onClick={handleClick} variant="outlined" color="primary" startIcon={<ArrowUpwardIcon />}>Thèmes</Button>;
};

const ThemeByDomShow = () => {
    const { permissions } = usePermissions();
    //const record = useRecordContext();
    const dataProvider = useDataProvider();
    const recordId = useGetRecordId();
    //console.log(recordId);
    const [dom] = useStore('dom', { select: {id:0 }, evaluation: {}, mask: {etu:false, nok:false} });
    const mask = dom.mask || { mask: {etu:false, nok:false} }; 
    const [state, setState] = useState({data: null, dom: 0, etu: mask.etu, nok: mask.nok});
    const [theme] = useTheme();


    useEffect(() => {
        //console.log('----ThemeByDomShow-------');
        //console.log('get: ' + MyConfig.API_URL + '/' + dom.select.id);
        dataProvider.getOne('themes/'+ recordId, {id: dom.select.id})
            .then(({ data }) => {
                //console.log('--dataProvider--');
                //console.log(data);
                setState({
                    data: data, etu: mask.etu, nok: mask.nok,
                });
            })
            .catch(error => {
                console.log('-error-', error);
                //window.location.href = MyConfig.BASE_PATH + '#/login';
                return;
            });

    },[dom]);

 

    if (!state.data) {
        return <div>...</div>;
    } else {
        let d = state.data;
        return (
            <>
                <Title title={d.name} /> 
                <ThemesButton /> 

                { (dom.select.id !== undefined && dom.select.id !== '0') ?
                    <Card>
                        <CardContent>
                            <h4>Périmètre {dom.select.name}</h4>
                            <RenderEval record={d.evaluation} theme={theme} />
                        </CardContent>
                    </Card>
                    : null }
                
                {permissions === 'admin' ?  <RegleMaskForm /> : null }

                <Card>
                    {permissions === 'admin' ?  <EditButton record={d} style={{float:'right'}}/> : null }
                    <CardContent>
                        <h2>{d.name}</h2>
                        <ReactMarkdown>{d.description !==''?d.description:d.descorig}</ReactMarkdown>
                        <ul>
                            {Array.isArray(d.regles) ? d.regles.map((r) => renderRegles(r, mask)) : null}
                        </ul>
                    </CardContent>
                </Card>
            </>
        );
    }

};

export default ThemeByDomShow;

