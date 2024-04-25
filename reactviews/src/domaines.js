import * as React from 'react';
import { List, Filter, Datagrid, TextField, ReferenceField, required } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput } from 'react-admin';
import { useNotify, useRedirect } from 'react-admin';


const DomaineFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Name" source="name" alwaysOn />
        <TextInput label="Description" source="description" alwaysOn />
        <ReferenceInput label="User" source="user" reference="users">
            <AutocompleteInput source="casid" optionText="casid" filterToQuery={searchText => ({ casid: searchText })} />
        </ReferenceInput>
        <ReferenceInput label="Parent" source="parent" reference="domaines">
            <AutocompleteInput source="name" optionText="name" filterToQuery={searchText => ({ name: searchText })} />
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

export const DomaineEdit =  () => {
    const parse = data => {
        return data?data:0;
    };
    const format = data => {
        return data?data:'';
    };
    return (
        <Edit mutationMode="pessimistic">
            <SimpleForm>
                <TextInput source="name" validate={required()} />
                <ReferenceInput label="Parent" source="parent" reference="domaines" parse={parse} format={format}>
                    <AutocompleteInput optionText="name" filterToQuery={searchText => ({ name: searchText })} />
                </ReferenceInput>
                <TextInput source="description" multiline fullWidth />
                <ReferenceInput label="User 1" source="user_1" reference="users" parse={parse} format={format}>
                    <AutocompleteInput source="casid" optionText="casid" filterToQuery={searchText => ({ casid: searchText })} />
                </ReferenceInput>
                <ReferenceInput label="User 2" source="user_2" reference="users" parse={parse} format={format}>
                    <AutocompleteInput optionText="casid" filterToQuery={searchText => ({ casid: searchText })}/>
                </ReferenceInput>
                <ReferenceInput label="User 3" source="user_3" reference="users" parse={parse} format={format}>
                    <AutocompleteInput optionText="casid" filterToQuery={searchText => ({ casid: searchText })} />
                </ReferenceInput>
            </SimpleForm>
        </Edit>
    );
};

export const DomaineCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    const parse = data => {
        return data?data:0;
    };
    const format = data => {
        return data?data:'';
    };
    const onSuccess = () => {
        notify('ra.notification.created', 'info', { undoable: false });
        redirect('list','domaines');
    };

    return (
        <Create mutationOptions={{ onSuccess }}>
            <SimpleForm>
                <TextInput source="name" validate={required()} />
                <ReferenceInput label="Parent" source="parent" parse={parse} format={format} reference="domaines">
                    <AutocompleteInput optionText="name" filterToQuery={searchText => ({ name: searchText })} />
                </ReferenceInput>
                <TextInput source="description" multiline fullWidth />
                <ReferenceInput label="User 1" source="user_1" reference="users" parse={parse} format={format}>
                    <AutocompleteInput optionText="casid" filterToQuery={searchText => ({ casid: searchText })} />
                </ReferenceInput>
                <ReferenceInput label="User 2" source="user_2" reference="users" parse={parse} format={format}>
                    <AutocompleteInput optionText="casid" filterToQuery={searchText => ({ casid: searchText })} />
                </ReferenceInput>
                <ReferenceInput label="User 3" source="user_3" reference="users">
                    <AutocompleteInput optionText="casid" filterToQuery={searchText => ({ casid: searchText })} />
                </ReferenceInput>
            </SimpleForm>
        </Create>
    );
};
