import * as React from 'react';
import { Card, CardContent } from '@mui/material';
import { Title } from 'react-admin';
import GitHubIcon from '@mui/icons-material/GitHub';
import Link from '@mui/material/Link';


export const About = () => (
    <Card>
        <Title title="About" />
        <CardContent>
            <h2>PSSImgmt</h2> 
            <div>Une application collaborative de GRC: <strong>G</strong>ouvernance, <strong>R</strong>isque, <strong>C</strong>onformité </div><br />
            <div style={{ display:'flex', justifyContent:'center' }}>
                <Link href="https://github.com/yvesago/pssimgmt" underline="hover" target="_blank">
                    <GitHubIcon/> https://github.com/yvesago/pssimgmt
                </Link>
            </div>
        </CardContent>
    </Card>
);


