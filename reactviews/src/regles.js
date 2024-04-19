import React, { Fragment, useState } from 'react';
import request from 'superagent';

import { List, Filter, Datagrid, TextField, SelectField, DateField, ReferenceField } from 'react-admin';
import { Create, Edit, SimpleForm, TextInput, NumberInput, Toolbar, SaveButton } from 'react-admin';
import { ReferenceArrayInput, SelectArrayInput, AutocompleteArrayInput, SelectInput, BooleanInput, FormDataConsumer } from 'react-admin';
import { useNotify, useRefresh, useRedirect, usePermissions, useRecordContext, useStore } from 'react-admin';
import { BulkDeleteButton, BulkExportButton, Labeled } from 'react-admin';
import { downloadCSV, useGetRecordId, useGetOne } from 'react-admin';

import jsonExport from 'jsonexport/dist';
import { NavLink } from 'react-router-dom';

import FactCheckIcon from '@mui/icons-material/FactCheck';
import AttachmentIcon from '@mui/icons-material/Attachment';
import Stack from '@mui/material/Stack';

import Divider from '@mui/material/Divider';
//import ListItemText from '@mui/material/ListItemText';
//import ListItem from '@mui/material/ListItem';
import Collapse from '@mui/material/Collapse';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';

import SliderInput from './SliderInput';

import ReactMarkdown from 'react-markdown';


import { MyConfig } from './MyConfig';

const classes = {
    field: {
        display: 'inline-block', width: '20em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
    buttons: {
        display: 'inline-block', float: 'left', margin: '4px' },
    docs: {
        display: 'inline-block', margin: '0em 0em 1em 0em' },
    ok: {
        textDecoration: 'none',
    },
    nok: {
        textDecoration: 'none',
        color: 'red'},
    etu: {
        textDecoration: 'none',
        color: 'green'},
    c0: {
        color: 'red'},
    c1: {
        color: 'orange'},
    c2: {
        color: 'orange'},
    c3: {
        color: 'green'},
};


const ShortTextField = (props) => {
    const record = useRecordContext();
    const c = {...classes[record.status],...classes.field};
    return ( <TextField sx={c} {...props} /> );
};

const status = [
    { name: 'ok', id: 'ok' },
    { name: 'non ok', id: 'nok' },
    { name: 'en étude', id: 'etu' },
];

const conform = [
    { name: 'Jamais', id: '0' },
    { name: 'Parfois', id: '1' },
    { name: 'Partiellement', id: '2' },
    { name: 'Totalement', id: '3' },
];

const conformChoice = (id) => {
    var res = ''; 
    conform.map((c) => {if (c.id === id) {res = c.name;} return res; });
    return res;
};

const axes = [
    { name: '', id: '0' },
    { name: 'Gouvernance', id: '1' },
    { name: 'Maîtrise des risques', id: '2' },
    { name: 'Maîtrise des systèmes', id: '3' },
    { name: 'Protection des systèmes', id: '4' },
    { name: 'Gestion des incidents', id: '5' },
    { name: 'Évaluation', id: '6' },
];

const axeChoice = (id) => {
    var res = ''; 
    axes.map((c) => {if (c.id === id) {res = c.name;} return res; });
    return res;
};

const ReglesFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Name" source="name" alwaysOn />
        <TextInput label="Description" source="descriptions" alwaysOn />
        <TextInput label="Code" source="code" />
        <TextInput label="Notes" source="notes" />
        <SelectInput source="status" choices={status} alwaysOn />
        <SelectInput source="axe" choices={axes} />
    </Filter>
);



const exporter = regles => {
    const reglesForExport = regles.map(regle => {
        // eslint-disable-next-line
        const { description, descorig, notes, regle_domaine, ordre, regles_iso, iso_ids, docs, docs_ids, theme, created_on, updated_on, created_by, updated_by, id, status, ...regleForExport } = regle; // omit entries
        regleForExport.descriptions = regle.description ? regle.description : regle.descorig; // add a field
        regleForExport.axe1= axeChoice(regle.axe1); // convert a field
        regleForExport.axe2= axeChoice(regle.axe2);
        if (regle.status === 'ok') {
            return regleForExport; }
        else { return null;}
    });
    for( var i = 0; i < reglesForExport.length; i++){ 
        if ( reglesForExport[i] === null ) { 
            reglesForExport.splice(i, 1); 
            i--;
        }
    }
    jsonExport(reglesForExport, {
        headers: ['code', 'name', 'descriptions'] // order fields in the export
    }, (err, csv) => {
        downloadCSV(csv, 'regles'); // download as 'posts.csv` file
    });
};

// eslint-disable-next-line
const ReglePanel = ({id, record, resource }) => (
    <Box>
        Status: <TextField source="status" /> <br />
	Notes: <br />
        <ShortTextField source="notes" record={record} /> <br />
        Axes : <SelectField source="axe1" choices={axes} />, <SelectField source="axe2" choices={axes} />
    </Box>
);

const RegleBulkActionButtons = () => {
    const { permissions } = usePermissions();
    return (
        <Fragment>
            <BulkExportButton />
            {permissions === 'admin' ?
                <BulkDeleteButton />
                : null}
        </Fragment>
    );
};


export const RegleList = () => (
    <List filters={<ReglesFilter />} exporter={exporter} bulkActionButtons={<RegleBulkActionButtons />} filterDefaultValues={{ status: 'ok' }} perPage={25} >
        <Datagrid expand={<ReglePanel />} rowClick="show">
            <TextField source="code" />
            <TextField source="name" />
            <ShortTextField source="descorig" />
            <ShortTextField source="description" />
            <DateField source="updated_on" />
        </Datagrid>
    </List>
); 

const RegleTitle = ({ record }) => {
    return <span>Règle {record ? `"${record.code}"` : ''}</span>;
};


export const RegleEdit = props => {
    const notify = useNotify();
    const refresh = useRefresh();
    const redirect = useRedirect();

    const onSuccess = () => {
        notify('Changes saved');
        redirect(`/regles/${props.id}/show`);
        refresh();
    };

    return (
        <Edit title={<RegleTitle />} {...props} onSubmit={onSuccess} mutationMode="pessimistic">
            <SimpleForm>
                <TextInput source="name" />
                <ReferenceField label="Périmètre" source="theme.id" reference="themes">
                    <TextField source="name" />
                </ReferenceField>
                <NumberInput source="ordre" />
                <TextInput source="code" />
                <TextInput multiline source="descorig" fullWidth />
                <TextInput multiline source="description" fullWidth />
                <TextInput multiline source="notes" fullWidth />
                <SelectInput source="status" choices={status} />
                <SelectInput source="axe2" choices={axes} />
                <SelectInput source="axe1" choices={axes} />

                <ReferenceArrayInput source="iso_ids" reference="isos" perPage={25}>
                    <AutocompleteArrayInput optionText="name" />
                </ReferenceArrayInput>

                <ReferenceArrayInput source="docs_ids" reference="docs" perPage={25}>
                    <SelectArrayInput optionText="name" />
                </ReferenceArrayInput>

                <Labeled label="Créé par">
                    <ReferenceField source="created_by" reference="users"><TextField source="name" /></ReferenceField>
                </Labeled>
                <Labeled label="Créé le">
                    <DateField source="created_on" />
                </Labeled>
                <Labeled label="Modifié par">
                    <ReferenceField source="updated_by" reference="users"><TextField source="name" /></ReferenceField>
                </Labeled>
                <Labeled label="Modifié le">
                    <DateField source="updated_on" />
                </Labeled>
            </SimpleForm>
        </Edit>
    );
};

export const RegleCreate = () => (
    <Create>
        <SimpleForm redirect="edit">
            <TextInput source="name" />
            <TextInput source="code" />
            <TextInput multiline source="descorig" fullWidth />
            <TextInput multiline source="description" fullWidth />
            <TextInput multiline source="notes" fullWidth />
            <SelectInput source="status" choices={status} />
            <SelectInput source="axe2" choices={axes} />
            <SelectInput source="axe1" choices={axes} />
            <TextInput source="regles_iso" />
            <ReferenceArrayInput source="docs_ids" reference="docs" perPage={100}>
                <SelectArrayInput optionText="name" />
            </ReferenceArrayInput>
        </SimpleForm>
    </Create>
);

const renderDoc = (doc, index) => (
    <span key={index}>
        {doc.url !== '' ?
            <a href={doc.url} style={{color: 'inherit', textDecoration:'none'}} rel="noreferrer" target="_blank">{doc.name}</a>
            : doc.name } 
    , </span> 
);

const renderIso = (iso, index) => (
    <Box key={index}>
        <span>
      [{iso.code}]  { iso.name }
        </span><br />
        {iso.descorig}
    </Box>
);



const ShowRegle = (r) => {
    var txt = '';
    if (r.description) {
        txt = r.description;
    }
    else {
        txt = r.descorig;
    }

    if ( r.regle_domaine && r.regle_domaine.modifdesc ) {
        txt = r.regle_domaine.modifdesc;
    }

    return (
        <NavLink to={`/regles/${r.id}/show`} style={{ color: 'inherit', textDecoration: 'none', }}>
            <div style={classes[r.status]} >
                <ReactMarkdown>{txt}</ReactMarkdown>
            </div>
        </NavLink>
    );
};

export const RegleView = (props) => {
    //console.log('**RegleView**');
    //console.log(props);
    const { permissions } = usePermissions();
    const record = props.record; //useRecordContext();
    const [dom] = useStore('dom', { select: {id:0 }, evaluation: {}, mask: {etu:false, nok:false} });
    if (record !== null ) {
        return (
            <>
                <Fragment>
                    <NavLink to={`/regles/${record.id}/show`} style={{ color: 'inherit', textDecoration: 'none', }}>
                        <strong>[{record.code}] { record.name }</strong>
                    </NavLink>
                </Fragment>
                <Typography variant="body2" component="span" display="block">
                    { ShowRegle(record) }
    
                    { record.regle_domaine 
                        ? ( record.regle_domaine.supldesc ? 
                            <Box><br />{record.regle_domaine.supldesc}</Box> : null )
                        : null
                    }
                    <div style={classes['docs']}>
                        <Stack alignItems="center" direction="row" gap={1}> 
                            { record.docs.length ? <><AttachmentIcon />Docs :</>  : null }
                            { Array.isArray(record.docs) ? record.docs.map((doc, index) => renderDoc(doc, index)): null }
                        </Stack>
                    </div>
                    <Divider variant="middle" />
                </Typography>
                { dom && dom.select.id && dom.select.id !== '0' ?
                    <Stack alignItems="center" direction="row" gap={1}> <FactCheckIcon color={(record.regle_domaine.applicable === 0)?'disabled':'primary'} />
                        { permissions === 'admin' ?
                            <ModifDialog record={record} dom={dom.select.id} viewonly={props.viewonly}/>
                            : null }
                        <EvalDialog  record={record} dom={dom.select.id} viewonly={props.viewonly}/>
                    </Stack>
                    : null
                }
                <RegleViewExt record={record} />
            </>
        );
    }
    return null;
};

//<EvalDialog  record={record} dom={dom.select.id} viewonly={props.viewonly} />
//<ModifDialog record={record} dom={dom.select.id} viewonly={props.viewonly} />

const Evaluation = (r) => {
    if (r) {
        const c = r.conform;
        const inherit = (r.applicable === 0)?'Non spécifique à ce périmètre':'';
        return (
            <>
                <Box sx={{ fontStyle: 'italic' }}>{inherit}</Box>
                Conforme : <span style={classes['c'+c]}>{conformChoice(c)}</span>, Évolution: {r.evolution}
            </>
        );
    }
    return( <Box sx={{ fontStyle: 'italic' }}>Non spécifique à ce périmètre</Box> );
};

const Updated_by  = (r) => {
    if (r.regle_domaine && r.regle_domaine.id) {
        return (
            <Typography variant="body2">
            Domain (eval) updated by user id : {r.regle_domaine.updated_by} on {new Date(r.regle_domaine.updated_on).toLocaleString()}
            </Typography>
        );
    } else {
        return (
            <Typography variant="body2">
            Updated by user id : {r.updated_by} on {new Date(r.updated_on).toLocaleString()}
            </Typography>
        );
    }
};

const RegleViewExt = (props) =>  {
    //console.log(props.record.id);
    const [open, setOpen] =useState(false);

    const handleClick = () => {
        setOpen(!open);
    };

    return (
        <Box>
            <Stack alignItems="center" direction="row" color="disabled" gap={1} onClick={handleClick}>
                {open ? <ExpandLess color="disabled" /> : <ExpandMore color="disabled" />}
                <Typography component="span" variant="caption">Détails</Typography>
            </Stack>
            <Collapse in={open} timeout="auto" key={10*props.record.id}>
                <Fragment key={props.record.id}>
                    {Updated_by(props.record)}
                    <Divider />
                    Direct URL: <NavLink to={`/regle/${props.record.code}`}>/regle/{props.record.code}</NavLink>
                    <Divider />
                    Axes: {axeChoice(props.record.axe1)}, {axeChoice(props.record.axe2)} 
                    <Divider />
                    Notes: <br />
                    {props.record.notes}

                    { props.record.regles_iso.length ? <span><Divider />Iso 27002: </span>  : null }
                    { Array.isArray(props.record.regles_iso) ? props.record.regles_iso.map((iso,index) => renderIso(iso,index)): null }
                </Fragment>
            </Collapse>
            <br />
        </Box>
    );
};

const RegleEditToolbar = () => (
    <Toolbar >
        <SaveButton />
    </Toolbar>
);

export const EvalDialog = (props) => {
    const [open, setOpen] = useState(false);
    const record = props.record;//useRecordContext();
    const refresh = useRefresh();
    const [rd, setRd] = useState(record.regle_domaine);
    const notify = useNotify();

    if (props.record.regle_domaine.id !== rd.id ) { 
        setRd(props.record.regle_domaine);
    }


    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleSave = (d) => {
        //console.log('===handleSave===');
        //console.log(d);
        const url = MyConfig.API_URL + '/regles/' + record.id + '/' + props.dom;
        var r = d;
        r.modif = 'eval';
        r.applicable = (d.applicable === true || d.applicable === 1 )?1:0;
        r.conform = r.conform.toString();
        r.evolution = r.evolution.toString();
        const token = localStorage.getItem('ttoken');
        request
            .put(url)
            .send(r)
            .set('Authorization', `Bearer ${token}`)
            .then( (response)  => {
                notify('Update');
                setRd(response.body);
                refresh();
            })
            .catch((e) => {
                console.log(e);
                notify('error', 'can\'t update');
            });

        setOpen(false);
    };

    if (props.record.status != 'ok') {
        return '';
    }

    if (props.viewonly === false) {
        return ( 
            <Typography component="span" variant="body2">
                <Evaluation {...rd} />
            </Typography>
        );
    }

    return (
        <>
            <Button variant="outlined" color="primary" onClick={handleClickOpen} sx={classes.buttons}>
        Évaluation
            </Button>
            <Evaluation {...rd} />
            <Dialog fullWidth={true} open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="dialog-title">Évaluation</DialogTitle>
                <DialogContent>

                    <Edit actions={null}>
                        <SimpleForm style={{marginLeft:'1em'}} toolbar={<RegleEditToolbar />} defaultValues={rd} onSubmit={handleSave}>

                            <BooleanInput label="Règle du périmètre" source="applicable" />

                            <FormDataConsumer>
                                {({ formData, ...rest }) => formData.applicable ? ( 
                                    <Box>
                                        <SliderInput label="Conforme" source="conform" choices={[
                                            { value: 0, label: 'Jamais'},
                                            { value: 1, label: 'Parfois' },
                                            { value: 2, label: 'Partiellement' },
                                            { value: 3, label: 'Totalement' },
                                        ]} {...rest}/>
                                        <br />

                                        <SliderInput label="Évolution" source="evolution" choices={[
                                            { value: -1, label: 'À arrêter'},
                                            { value: 0, label: 'Fix' },
                                            { value: 1, label: 'À faire' },
                                        ]}  {...rest} />

                                    </Box>
                                ) : null
                                }
                            </FormDataConsumer>

                        </SimpleForm>
                    </Edit>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color="primary">Cancel</Button>
                </DialogActions>
            </Dialog>
        </>
    );

};

export const ModifDialog = (props) => {
    const [open, setOpen] = useState(false);
    const notify = useNotify();
    const refresh = useRefresh();


    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleSave = (d) => {
        const url = MyConfig.API_URL + '/regles/' + d.regle + '/' + d.domaine_id;
        var r = d;
        r.modif = 'modif';
        const token = localStorage.getItem('ttoken');
        request
            .put(url)
            .send(r)
            .set('Authorization', `Bearer ${token}`)
            .then( ()  => {
                //setRd(response.body);
                refresh();
                notify('Update');
            })
            .catch((e) => {
                console.log(e);
                notify('error', 'can\'t update');
            });

        setOpen(false);
    };

    if (props.viewonly === false) {
        return null;
    }

    return (
        <Typography component="span" variant="body2">
            <Button variant="outlined" color="primary" onClick={handleClickOpen} sx={classes.buttons}>
        Modification
            </Button>
            <Dialog fullWidth={true} open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
                <DialogContent>

                    <Edit actions={null}>
                        <SimpleForm toolbar={<RegleEditToolbar />} defaultValues={props.record.regle_domaine} onSubmit={handleSave}>
                            <TextInput source="modifdesc" multiline fullWidth />
                            <TextInput source="supldesc" multiline fullWidth />
                        </SimpleForm>
                    </Edit>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color="primary">Cancel</Button>
                </DialogActions>
            </Dialog>
        </Typography>
    );

};

export const RegleCodeShow = () => {
    const recordId = useGetRecordId();
    const redirect = useRedirect();
    const { data: data, isLoading, error } = useGetOne('regle', { id: recordId });
    if (error) return <Error />;
    if (!data) return null;
    redirect('show','regles', data.id);
};
