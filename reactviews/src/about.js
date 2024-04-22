import * as React from 'react';
import { Card, CardContent } from '@mui/material';
import { Title } from 'react-admin';
import GitHubIcon from '@mui/icons-material/GitHub';
import Link from '@mui/material/Link';
import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import TableHead from '@mui/material/TableHead';
import TableBody from '@mui/material/TableBody';
//import Paper from '@mui/material/Paper';

const styles = {
    field: {
        display: 'inline-block', width: '20em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
    c0: {
        color: 'red'},
    c1: {
        color: 'orange'},
    c2: {
        color: 'orange'},
    c3: {
        color: 'green'},
};

const conform = [
    { name: 'Totalement', id: '3', appliqué: '✔', strictement: '✔', documenté :'✔', contrôlée :'✔' },
    { name: 'Partiellement', id: '2', appliqué: '✔', strictement: '✔', documenté :'~', contrôlée :'.' },
    { name: 'Parfois', id: '1', appliqué: '✔', strictement: '.', documenté :'.', contrôlée :'.' },
    { name: 'Jamais', id: '0', appliqué: '.', strictement: '.', documenté :'.', contrôlée :'.' },
];



export const About = () => {
    return (<Card>
        <Title title="About" />
        <CardContent>
            <h2>PSSImgmt</h2> 
            <div>Une application collaborative de GRC: <strong>G</strong>ouvernance, <strong>R</strong>isque, <strong>C</strong>onformité </div><br />
            <h3>Conformité des régles</h3>
            <TableContainer>
                <Table sx={{ width: 250 }} aria-label="conform table">
                    <TableHead>
                        <TableCell><strong>Conformité</strong></TableCell>
                        <TableCell>Appliquée</TableCell>
                        <TableCell>Systématiquement</TableCell>
                        <TableCell>Documentée</TableCell>
                        <TableCell>Contrôlée</TableCell>
                    </TableHead>
                    <TableBody>{conform.map((row) => (
                        <TableRow>
                            {Object.keys(row).map((k,i)=>{
                                if (k === 'id') return;
                                return(<TableCell align="left" sx={styles['c'+row.id]}>{row[k]}</TableCell>);
                            })}
                        </TableRow>
                    ))}</TableBody>
                </Table>
            </TableContainer>
            <h3>Code source v2.x</h3>
            <div style={{ display:'flex', justifyContent:'center' }}>
                <Link href="https://github.com/yvesago/pssimgmt" underline="hover" target="_blank">
                    <GitHubIcon/> https://github.com/yvesago/pssimgmt
                </Link>
            </div>
        </CardContent>
    </Card>);
};


