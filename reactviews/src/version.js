import * as React from 'react';
import { Fragment } from 'react';
import { List, Filter, Datagrid, TextField, DateField, ReferenceField } from 'react-admin';
import { Edit, Create, SimpleForm, TextInput, DateInput, Labeled, required } from 'react-admin';
import { BulkExportButton, BulkDeleteButton } from 'react-admin';
import { useNotify, useRedirect, usePermissions } from 'react-admin';

const VersionFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Nom" source="name" alwaysOn />
        <TextInput label="Changements" source="changelog" />
    </Filter>
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

export const VersionList = () => (
    <List filters={<VersionFilter />} bulkActionButtons={<RegleBulkActionButtons />} sort={{ field: 'created_on', order: 'DESC' }}>
        <Datagrid rowClick="edit">
            <TextField source="name" label="Nom" />
            <DateField source="validationdate"  label="Date de validation" />
            <TextField source="validationpar" label="Valideurs" />
            <TextField source="status" />
            <TextField source="changelog" label="Changements" />
            <DateField source="created_on"  label="Créé" />
            <DateField source="updated_on"  label="Modifié" />
        </Datagrid>
    </List>
);

const dateParseRegex = /(\d{4})-(\d{2})-(\d{2})/;
const dateFormatRegex = /^\d{4}-\d{2}-\d{2}$/;

const convertDateToString = (value) => {
    // value is a `Date` object
    if (!(value instanceof Date) || isNaN(value.getDate())) return '';
    const pad = '00';
    const yyyy = value.getFullYear().toString();
    const MM = (value.getMonth() + 1).toString();
    const dd = value.getDate().toString();
    return `${yyyy}-${(pad + MM).slice(-2)}-${(pad + dd).slice(-2)}`;
};

const dateFormatter = (value) => {
    // null, undefined and empty string values should not go through dateFormatter
    // otherwise, it returns undefined and will make the input an uncontrolled one.
    if (value == null || value === '') return '';
    if (value instanceof Date) return convertDateToString(value);
    // Valid dates should not be converted
    if (dateFormatRegex.test(value)) return value;

    return convertDateToString(new Date(value));
};

const dateParser = value => {
    //console.log(value);
    //value is a string of "YYYY-MM-DD" format
    const match = dateParseRegex.exec(value);
    if (match === null || match.length === 0) return;
    //console.log(match);
    const d = new Date(parseInt(match[1]), parseInt(match[2], 10) - 1, parseInt(match[3]));
    //console.log(d);
    if (isNaN(d.getDate())) return;
    return d;
};

export const VersionEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="name" label="Nom" validate={required()} />
            <TextInput source="validationpar" label="Validation par" fullWidth />
            <DateInput source="validationdate" label="Date de validation" parse={dateParser} format={dateFormatter} />
            <TextInput multiline source="changelog" label="Changements" fullWidth />
            <TextInput source="status" defaultValue="redaction" />
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

export const VersionCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    const onSuccess = () => {
        notify('ra.notification.created', { type: 'info', undoable: false });
        redirect('/versions/');
    };

    return (
        <Create mutationOptions={{onSuccess}}>
            <SimpleForm>
                <TextInput source="name" validate={required()} />
                <TextInput source="validationpar" fullWidth />
                <DateInput source="validationdate" parse={dateParser} />
                <TextInput multiline source="changelog" fullWidth />
            </SimpleForm>
        </Create>
    );
};
