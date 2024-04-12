import React from 'react';

import RegleByDomShow from './reglesByDomShow';
import MenuDomaine from './MenuDomaine';
import { Card, Grid, Box } from '@mui/material';

const RegleShow = () => {
    return (
        <Box sx={{ flexGrow: 1 }} >
            <Grid container spacing={2}>
                <Grid item xs={3}>
                    <Card><MenuDomaine /></Card>
                </Grid>
                <Grid item xs={9}>
                    <Card><RegleByDomShow /></Card>
                </Grid>
            </Grid>
        </Box>
    );
};

export default RegleShow;

