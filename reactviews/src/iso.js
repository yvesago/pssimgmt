import * as React from 'react';
import { List, Filter, Datagrid, TextField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput } from 'react-admin';
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

export const IsoEdit = props => (
    <Edit {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <TextInput source="code" />
            <TextInput multiline source="descorig" fullWidth />
        </SimpleForm>
    </Edit>
);

export const IsoCreate = props => {
    const notify = useNotify();
    const redirect = useRedirect();

    const onSuccess = (data) => {
        notify('ra.notification.created', 'info', { smart_count: 1 }, props.mutationMode === 'undoable');
        redirect('list', props.basePath, data.id, data);
    };

    return (
        <Create onSuccess={onSuccess} {...props}>
            <SimpleForm>
                <TextInput source="name" />
                <TextInput source="code" />
                <TextInput multiline source="descorig" fullWidth />
            </SimpleForm>
        </Create>
    );
};
