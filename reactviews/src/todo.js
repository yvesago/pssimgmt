import * as React from 'react';
import { List, Filter, Datagrid, TextField, NumberField, DateField, ReferenceField } from 'react-admin';
//import { Edit, Create, SimpleForm, TextInput, SelectInput, NumberInput, Labeled } from 'react-admin';
import { TextInput, FunctionField, SelectInput, ReferenceInput, AutocompleteInput  } from 'react-admin';
import { WrapperField, BulkDeleteButton } from 'react-admin';
import { useNotify, useRedirect, useRecordContext, usePermissions } from 'react-admin';


import Box from '@mui/material/Box';

/*const status = [
    { name: 'ok', id: 'ok' },
    { name: 'non ok', id: 'nok' },
    { name: 'en étude', id: 'etu' },
];*/


const evol = [
    { id: '1', name: 'À faire' },
    { id: '0', name: 'Fix' },
    { id: '-1', name: 'À arrêter' },
];

const evolChoice = (id) => {
    var res = '';
    evol.map((c) => {if (c.id === id) {res = c.name;} return res; });
    return res;
};

const EvolTextField = () => {
    const record = useRecordContext();
    //console.log(record);
    if (record && record.id === 0 ) { return (''); }
    return ( <Box>{evolChoice(record.evolution)}</Box> );
};

const intApplicable = [
    { id: '1', name: '✔' },
    { id: '0', name: '✘' },
];

const filterToQuery = searchText => ({ name: `${searchText}` });
const filterCSSIQuery = searchText => ({ name: `${searchText}` });

const TodoFilter = (props) => {
    const { permissions } = usePermissions();
    return (<Filter {...props}>
        <SelectInput label="Applicable" source="applicable" choices={intApplicable} alwaysOn />
        <SelectInput label="Conforme" source="conform" choices={conform} alwaysOn />
        <SelectInput label="Évolution" source="evolution" choices={evol} alwaysOn />
        <ReferenceInput source="domaine_id" reference="domaines" label="Périmètre">
            <AutocompleteInput label="Périmètre" optionText="name" filterToQuery={filterToQuery} />
        </ReferenceInput>
        {permissions === 'admin' ?
            <ReferenceInput source="domaine_id" reference="domaines" label="CSSIs">
                <ReferenceInput source="user" reference="users">
                    <AutocompleteInput label="CSSIs" optionText="name" filterToQuery={filterCSSIQuery} />
                </ReferenceInput>
            </ReferenceInput>
            : null}
    </Filter>);
};

const styles = {
    field: {
        display: 'inline-block', width: '20em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
    c0: {
        color: 'red'},
    c1: {
        color: 'orange'},
    c2: {
        color: 'orange'},
    c3: {
        color: 'green'},
};

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


const ConformTextField = () => {
    const record = useRecordContext();
    //console.log(record);
    if (record && record.id === 0 ) { return (''); }
    return ( <Box sx={styles['c'+record.conform]}>{conformChoice(record.conform)}</Box> );
};


const ShortTextField = (props) => {
    return ( <TextField sx={ styles.field } {...props} /> );
};


const ModifPanel = () => {
    const record = useRecordContext();
    //console.log(record);
    return (
        <Box>
            Updated on <DateField source="updated_on" record={record} /> by &nbsp; 
            <ReferenceField source="updated_by" reference="users" link={false}>
                <TextField source="name" />
            </ReferenceField> <br />
            <ShortTextField source="modifdesc" record={record} /> <br />
            <ShortTextField source="supldesc" record={record} />
        </Box>
    );
};


export const TodoList = () => {
    const { permissions } = usePermissions();
    return (<List
        filters={<TodoFilter />}
        filterDefaultValues={{ applicable: 1, evolution: '1' }}
        sort={{ field: 'conform', order: 'ASC' }}
        exporter={false}
        perPage={50}
    >
        <Datagrid expand={<ModifPanel />} bulkActionButtons={permissions === 'admin' ? <BulkDeleteButton /> : false}>
            <FunctionField
                label='App.'
                source="applicable"
                render={record => (record.applicable === 1)?'✔':'✘'}
            />
            <ReferenceField label="Périmètre" source="domaine_id" reference="domaines" link={permissions === 'admin' ? 'edit' : false}>
                <TextField source="name" />
            </ReferenceField>
            <ReferenceField label="Règle" source="regle" reference="regles" link="show">
                <TextField source="name" />
            </ReferenceField>
            <WrapperField label="CSSIs" sortBy="user_1">
                <ReferenceField label="CSSI 1" source="user_1" reference="users" link={permissions === 'admin' ? 'edit' : 'show'}>
                    <TextField source="name" />
                </ReferenceField> &nbsp;
                <ReferenceField label="CSSI 2" source="user_2" reference="users" link={permissions === 'admin' ? 'edit' : 'show'}>
                    <>, &nbsp; <TextField source="name" /></>
                </ReferenceField>
                <ReferenceField label="CSSI 3" source="user_3" reference="users" link={permissions === 'admin' ? 'edit' : 'show'}>
                    <>, &nbsp; <TextField source="name" /></>
                </ReferenceField>
            </WrapperField>
            <ConformTextField source="conform" />
            <EvolTextField source="evolution" />
        </Datagrid>
    </List>);
};

