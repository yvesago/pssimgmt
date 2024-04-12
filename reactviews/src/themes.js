import * as React from 'react';
import { List, Filter, Datagrid, TextField, NumberField, DateField, ReferenceField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, SelectInput, NumberInput, Labeled } from 'react-admin';
import { ReferenceInput, AutocompleteInput, ReferenceArrayInput, AutocompleteArrayInput } from 'react-admin';
import { useNotify, useRedirect } from 'react-admin';


const status = [
    { name: 'ok', id: 'ok' },
    { name: 'non ok', id: 'nok' },
    { name: 'en étude', id: 'etu' },
];

const ThemeFilter = (props) => (
    <Filter {...props}>
        <TextInput label="name" source="name" alwaysOn />
        <TextInput label="description" source="descriptions" />
    </Filter>
);

const styles = {
    field: {
        display: 'inline-block', width: '20em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
};

const ShortTextField = (props) => {
    return ( <TextField sx={ styles.field } {...props} /> );
};


import { useRecordContext } from 'react-admin';

const FixEmptyField = ({ source }) => {
    const record = useRecordContext();
    if (record && record.id === 0 ) { return (''); }
    return (<span>{record && record[source]}</span>);
};


export const ThemeList = () => (
    <List filters={<ThemeFilter />} perPage={25}>
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <NumberField source="ordre" />
            <ReferenceField label="Parent" source="parent" reference="themes">
                <FixEmptyField source="name" />
            </ReferenceField>
            <ShortTextField source="descorig" />
            <ShortTextField source="description" />
            <TextField source="notes" />
            <TextField source="status" />
            <DateField source="updated_on" />
        </Datagrid>
    </List>
);

export const ThemeEdit = props => {
    const parse = data => {
        return data?data:0;
    };
    const format = data => {
        return data?data:'';
    };
    return (<Edit {...props}>
        <SimpleForm>
            <TextField source="id" />
            <TextInput source="name" />
            <NumberInput source="ordre" />
            <ReferenceInput label="Parent" source="parent" parse={parse} format={format} reference="themes">
                <AutocompleteInput optionText="name" />
            </ReferenceInput>

            <TextInput multiline source="descorig" fullWidth />
            <TextInput multiline source="description" fullWidth />
            <TextInput multiline source="notes" fullWidth />
            <SelectInput source="status" choices={status} />

            <ReferenceArrayInput source="regles_ids" reference="regles">
                <AutocompleteArrayInput optionText="code" />
            </ReferenceArrayInput>

            <ReferenceArrayInput source="iso_ids" reference="isos" perPage={25} filterToQuery={searchText => ({ name: searchText })}>
                <AutocompleteArrayInput optionText="name" />
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
    );};

export const ThemeCreate = props => {
    const notify = useNotify();
    const redirect = useRedirect();

    const parse = data => {
        return data?data:0;
    };
    const format = data => {
        return data?data:'';
    };
    const onSuccess = (data) => {
        notify('ra.notification.created', 'info', { smart_count: 1 }, props.mutationMode === 'undoable');
        redirect('list', props.basePath, data.id, data);
    };

    return (
        <Create onSuccess={onSuccess} {...props}>
            <SimpleForm>
                <TextInput source="name" />
                <NumberInput source="ordre" />
                <ReferenceInput label="Parent" source="parent" parse={parse} format={format} reference="themes">
                    <AutocompleteInput optionText="name" />
                </ReferenceInput>
                <TextInput multiline source="descorig" fullWidth />
                <TextInput multiline source="description" fullWidth />
                <TextInput multiline source="notes" fullWidth />
                <SelectInput source="status" choices={status} />
            </SimpleForm>
        </Create>
    );

};
