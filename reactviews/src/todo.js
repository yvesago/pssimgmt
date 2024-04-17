import * as React from 'react';
import { List, Filter, Datagrid, TextField, NumberField, DateField, ReferenceField } from 'react-admin';
//import { Edit, Create, SimpleForm, TextInput, SelectInput, NumberInput, Labeled } from 'react-admin';
import { TextInput, NumberInput, SelectInput, ReferenceInput, AutocompleteInput  } from 'react-admin';
//import { ReferenceInput, AutocompleteInput, ReferenceArrayInput, AutocompleteArrayInput } from 'react-admin';
import { useNotify, useRedirect } from 'react-admin';


/*const status = [
    { name: 'ok', id: 'ok' },
    { name: 'non ok', id: 'nok' },
    { name: 'en étude', id: 'etu' },
];*/

const filterToQuery = searchText => ({ name: `${searchText}` });

const TodoFilter = (props) => (
    <Filter {...props}>
        <SelectInput label="Applicable" source="applicable" choices={[
            { id: '1', name: '1' },
            { id: '0', name: '0' },
        ]} alwaysOn />
        <ReferenceInput source="domaine_id" reference="domaines" alwaysOn>
            <AutocompleteInput label="Périmètre" optionText="name" filterToQuery={filterToQuery} />
        </ReferenceInput>
    </Filter>
);

const styles = {
    field: {
        display: 'inline-block', width: '20em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
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


export const TodoList = () => (
    <List filters={<TodoFilter />} filterDefaultValues={{ applicable: 1 }} perPage={50}>
        <Datagrid>
            <ReferenceField label="Périmètre" source="domaine_id" reference="domaines">
                <TextField source="name" />
            </ReferenceField>
            <ReferenceField label="Règle" source="regle" reference="regles" link="show">
                <TextField source="name" />
            </ReferenceField>
            <NumberField source="applicable" />
            <TextField source="conform" />
            <TextField source="evolution" />
            <ShortTextField source="modifdesc" />
            <ShortTextField source="supldesc" />
            <DateField source="updated_on" />
        </Datagrid>
    </List>
);

