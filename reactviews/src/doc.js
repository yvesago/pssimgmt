import * as React from 'react';
import { List, Filter, Datagrid, TextField, DateField, NumberField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, NumberInput } from 'react-admin';
import { useNotify, useRedirect } from 'react-admin';

const DocFilter = (props) => (
    <Filter {...props}>
        <TextInput label="name" source="name" alwaysOn />
        <TextInput label="description" source="description" />
    </Filter>
);



export const DocList = () => (
    <List filters={<DocFilter />} >
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <TextField source="url" />
            <TextField source="description" />
            <TextField source="notes" />
            <TextField source="status" />
        </Datagrid>
    </List>
);

export const DocEdit = props => (
    <Edit {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <NumberInput source="ordre" />
            <TextInput source="url" fullWidth />
            <TextInput multiline source="description" fullWidth />
            <TextInput multiline source="notes" fullWidth />
            <TextInput source="status" />
            <NumberField source="created_by" />
            <DateField source="created_on" />
            <NumberField source="updated_by" />
            <DateField source="updated_on" />
        </SimpleForm>
    </Edit>
);

export const DocCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    const onSuccess = (data) => {
        notify('ra.notification.created', { type: 'info', undoable: true });
        redirect(`/docs/${data.id}`);
    };

    return (
        <Create mutationOptions={{onSuccess}}>
            <SimpleForm>
                <TextInput source="name" />
                <TextInput source="url" fullWidth />
                <TextInput multiline source="description" fullWidth />
                <TextInput multiline source="notes" fullWidth />
            </SimpleForm>
        </Create>
    );
};
