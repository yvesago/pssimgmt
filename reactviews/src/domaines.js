import * as React from 'react';
import { List, Filter, Datagrid, TextField, ReferenceField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, ReferenceInput, SelectInput, AutocompleteInput } from 'react-admin';
import { useNotify, useRedirect } from 'react-admin';


const DomaineFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Name" source="name" alwaysOn />
        <TextInput label="Description" source="description" alwaysOn />
        <ReferenceInput label="User" source="user" reference="users" filterToQuery={searchText => ({ casid: searchText })}>
            <AutocompleteInput source="casid" optionText="casid" />
        </ReferenceInput>
        <ReferenceInput label="Parent" source="parent" reference="domaines" filterToQuery={searchText => ({ name: searchText })}>
            <AutocompleteInput source="name" optionText="name" />
        </ReferenceInput>
    </Filter>
);

import { useRecordContext } from 'react-admin';

const FixEmptyField = ({ source }) => {
    const record = useRecordContext();
    if (record && record.id === 0 ) { return (''); }
    return (<span>{record && record[source]}</span>);
};


export const DomaineList = () => (
    <List filters={<DomaineFilter />} perPage={25}>
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <ReferenceField label="Parent" source="parent" reference="domaines" >
                <FixEmptyField source="name" />
            </ReferenceField>
            <TextField source="description" />
            <ReferenceField label="User 1" source="user_1" reference="users">
                <TextField source="casid" />
            </ReferenceField>
            <ReferenceField label="User 2" source="user_2" reference="users">
                <TextField source="casid" />
            </ReferenceField>
            <ReferenceField label="User 3" source="user_3" reference="users">
                <TextField source="casid" />
            </ReferenceField>
        </Datagrid>
    </List>
);

export const DomaineEdit =  (props) => {
    const parse = data => {
        return data?data:0;
    };
    const format = data => {
        return data?data:'';
    };
    return (
        <Edit mutationMode="pessimistic" {...props}>
            <SimpleForm>
                <TextInput source="name" />
                <ReferenceInput label="Parent" source="parent" reference="domaines" parse={parse} format={format} filterToQuery={searchText => ({ name: searchText })}>
                    <SelectInput optionText="name" />
                </ReferenceInput>
                <TextInput source="description" multiline fullWidth />
                <ReferenceInput label="User 1" source="user_1" reference="users" parse={parse} format={format} filterToQuery={searchText => ({ casid: searchText })}>
                    <AutocompleteInput source="casid" optionText="casid" />
                </ReferenceInput>
                <ReferenceInput label="User 2" source="user_2" reference="users" parse={parse} format={format} filterToQuery={searchText => ({ casid: searchText })}>
                    <AutocompleteInput optionText="casid" />
                </ReferenceInput>
                <ReferenceInput label="User 3" source="user_3" reference="users" parse={parse} format={format} filterToQuery={searchText => ({ casid: searchText })} >
                    <AutocompleteInput optionText="casid" />
                </ReferenceInput>
            </SimpleForm>
        </Edit>
    );
};

export const DomaineCreate = props => {
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
                <ReferenceInput label="Parent" source="parent" parse={parse} format={format} reference="domaines">
                    <SelectInput optionText="name" />
                </ReferenceInput>
                <TextInput source="description" multiline fullWidth />
                <ReferenceInput label="User 1" source="user_1" reference="users" parse={parse} format={format} filterToQuery={searchText => ({ casid: searchText })}>
                    <AutocompleteInput optionText="casid" />
                </ReferenceInput>
                <ReferenceInput label="User 2" source="user_2" reference="users" parse={parse} format={format} filterToQuery={searchText => ({ casid: searchText })}>
                    <AutocompleteInput optionText="casid" />
                </ReferenceInput>
                <ReferenceInput label="User 3" source="user_3" reference="users" filterToQuery={searchText => ({ casid: searchText })}>
                    <AutocompleteInput optionText="casid" />
                </ReferenceInput>
            </SimpleForm>
        </Create>
    );
};