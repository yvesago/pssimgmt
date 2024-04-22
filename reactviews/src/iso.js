import * as React from 'react';
import { List, Filter, Datagrid, TextField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, required } from 'react-admin';
import { useNotify, useRedirect } from 'react-admin';

const IsoFilter = (props) => (
    <Filter {...props}>
        <TextInput label="name" source="name" alwaysOn />
        <TextInput label="description" source="descorig" />
    </Filter>
);



export const IsoList = props => (
    <List filters={<IsoFilter />} perPage={25} {...props}>
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <TextField source="code" />
            <TextField source="descorig" />
        </Datagrid>
    </List>
);

export const IsoEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="name" validate={required()} />
            <TextInput source="code" />
            <TextInput multiline source="descorig" fullWidth />
        </SimpleForm>
    </Edit>
);

export const IsoCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    const onSuccess = () => {
        notify('ra.notification.created', { type: 'info', undoable: false });
        redirect('/isos');
    };

    return (
        <Create mutationOptions={{onSuccess}}>
            <SimpleForm>
                <TextInput source="name" validate={required()} />
                <TextInput source="code" />
                <TextInput multiline source="descorig" fullWidth />
            </SimpleForm>
        </Create>
    );
};
