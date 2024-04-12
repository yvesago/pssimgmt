import * as React from 'react';
import { List, Filter, Datagrid, TextField, DateField, ReferenceField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, Labeled, Show, SimpleShowLayout } from 'react-admin';
import { useNotify, useRedirect, useRecordContext  } from 'react-admin';

import ReactMarkdown from 'react-markdown';


const ShortTextField = (props) => {
    const c = {display: 'inline-block', width: '20em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'};
    return ( <TextField sx={c} {...props} /> );
};

const DescShow = () => {
    const record = useRecordContext();
    if (!record) return null;
    return ( <ReactMarkdown>{record.description}</ReactMarkdown> );
};

const NotesShow = () => {
    const record = useRecordContext();
    if (!record) return null;
    return ( <ReactMarkdown>{record.notes}</ReactMarkdown> );
};

const DocumentFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Nom" source="name" alwaysOn />
        <TextInput label="Changements" source="changelog" />
    </Filter>
);



export const DocumentList = () => (
    <List filters={<DocumentFilter />} sort={{ field: 'created_on', order: 'DESC' }}>
        <Datagrid rowClick="show">
            <TextField source="name" label="Nom" />
            <TextField source="titre"  label="Titre" />
            <ShortTextField source="description" label="Description" />
            <TextField source="notes" />
            <DateField source="created_on"  label="Créé" />
            <DateField source="updated_on"  label="Modifié" />
        </Datagrid>
    </List>
);


export const DocumentEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="name" label="Nom" />
            <TextInput source="titre" label="Titre" />
            <TextInput multiline source="description" label="Description" fullWidth />
            <TextInput multiline source="notes" label="Notes" fullWidth />
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

export const DocumentCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    const onSuccess = () => {
        notify('ra.notification.created', { type: 'info', undoable: false });
        redirect('/documents/');
    };

    return (
        <Create mutationOptions={{onSuccess}}>
            <SimpleForm>
                <TextInput source="name" />
                <TextInput source="titre" fullWidth />
                <TextInput multiline source="description" fullWidth />
                <TextInput multiline source="notes" fullWidth />
            </SimpleForm>
        </Create>
    );
};

export const DocumentShow = () => (
    <Show>
        <SimpleShowLayout>
            <Labeled label="Nom">
                <TextField source="name" />
            </Labeled>
            <Labeled label="Titre">
                <TextField source="titre" />
            </Labeled>
            <Labeled label="Description">
                <DescShow />
            </Labeled>
            <Labeled label="Notes">
                <NotesShow />
            </Labeled>
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
        </SimpleShowLayout>
    </Show>
);
