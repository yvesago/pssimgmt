import * as React from 'react';
import { List, Filter, Datagrid, TextField, EmailField } from 'react-admin';
import { ReferenceArrayField, SingleFieldList, ChipField } from 'react-admin';
import { Edit, SimpleForm, TextInput, SelectInput  } from 'react-admin';

const roles = [
    { name: 'admin', id: 'admin' },
    { name: 'cssi', id: 'cssi' },
    { name: 'reader', id: 'reader' },
    { name: 'guest', id: 'guest' },
];

const UserFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Login" source="casid" alwaysOn />
        <SelectInput source="user_role" choices={roles} alwaysOn />
        <TextInput label="Nom" source="name"  />
    </Filter>
);

export const UserList = () => (
    <List filters={<UserFilter />} perPage={25}>
        <Datagrid rowClick="edit">
            <TextField source="casid" />
            <TextField source="name" />
            <TextField source="user_role" />
            <EmailField source="email" />
        </Datagrid>
    </List>
);

export const UserEdit = () => (
    <Edit>
        <SimpleForm>
            <TextField label="Login" source="casid" />
            <TextInput source="name" />
            <SelectInput  source="user_role" choices={roles} />
            <TextInput source="email" />
            <ReferenceArrayField label="Domaines" reference="domaines" source="doms">
                <SingleFieldList>
                    <ChipField source="name" />
                </SingleFieldList>
            </ReferenceArrayField>
        </SimpleForm>
    </Edit>
);

