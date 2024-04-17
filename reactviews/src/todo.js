import * as React from 'react';
import { List, Filter, Datagrid, TextField, NumberField, DateField, ReferenceField } from 'react-admin';
//import { Edit, Create, SimpleForm, TextInput, SelectInput, NumberInput, Labeled } from 'react-admin';
import { TextInput, NumberInput, SelectInput, ReferenceInput, AutocompleteInput  } from 'react-admin';
import { BulkExportButton, BulkDeleteButton } from 'react-admin';
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

const filterToQuery = searchText => ({ name: `${searchText}` });
const filterCSSIQuery = searchText => ({ name: `${searchText}` });

const TodoFilter = (props) => {
    //const { permissions } = usePermissions();
    return (<Filter {...props}>
        <SelectInput label="Applicable" source="applicable" choices={[
            { id: '1', name: '1' },
            { id: '0', name: '0' },
        ]} alwaysOn />
        <SelectInput label="Évolution" source="evolution" choices={evol} alwaysOn />
        <ReferenceInput source="domaine_id" reference="domaines" alwaysOn>
            <AutocompleteInput label="Périmètre" optionText="name" filterToQuery={filterToQuery} />
        </ReferenceInput>
        <ReferenceInput source="domaine_id" reference="domaines" alwaysOn>
            <ReferenceInput source="user_1" reference="users">
                <AutocompleteInput label="CSSI 1" optionText="name" filterToQuery={filterCSSIQuery} />
            </ReferenceInput>
        </ReferenceInput>
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

/*import { useRecordContext } from 'react-admin';

const FixEmptyField = ({ source }) => {
    const record = useRecordContext();
    if (record && record.id === 0 ) { return (''); }
    return (<span>{record && record[source]}</span>);
};*/


const ModifPanel = () => {
    const record = useRecordContext();
    //console.log(record);
    return (
        <Box>
            Updated on <DateField source="updated_on" record={record} /> by <TextField source="updated_by" record={record} /> <br />
            <ShortTextField source="modifdesc" record={record} /> <br />
            <ShortTextField source="supldesc" record={record} />
        </Box>
    );
};


const TodoBulkActionButtons = () => {
    const { permissions } = usePermissions();
    return (
        <Box>
            <BulkExportButton />
            {permissions === 'admin' ?
                <BulkDeleteButton />
                : null}
        </Box>
    );
};


export const TodoList = () => {
    //const { permissions } = usePermissions();
    return (<List
        filters={<TodoFilter />}
        filterDefaultValues={{ applicable: 1, evolution: '1' }}
        sort={{ field: 'conform', order: 'ASC' }}
        bulkActionButtons={<TodoBulkActionButtons />}
        perPage={50}
    >
        <Datagrid expand={<ModifPanel />}>
            <ReferenceField label="Périmètre" source="domaine_id" reference="domaines">
                <TextField source="name" />
            </ReferenceField>
            <ReferenceField label="Règle" source="regle" reference="regles" link="show">
                <TextField source="name" />
            </ReferenceField>
            <ReferenceField label="CSSI 1" source="user_1" reference="users">
                <TextField source="name" />
            </ReferenceField>
            <NumberField source="applicable" />
            <ConformTextField source="conform" />
            <EvolTextField source="evolution" />
        </Datagrid>
    </List>);
};

