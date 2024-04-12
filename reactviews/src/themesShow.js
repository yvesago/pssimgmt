import React from 'react';

import ThemeByDomShow from './themesByDomShow';
import MenuDomaine from './MenuDomaine';
import { Card, Grid, Box } from '@mui/material';

const ThemeShow = () => {
    return (
        <Box sx={{ flexGrow: 1 }} >
            <Grid container spacing={2}>
                <Grid item xs={3}>
                    <Card><MenuDomaine /></Card>
                </Grid>
                <Grid item xs={9}>
                    <ThemeByDomShow />
                </Grid>
            </Grid>
        </Box>
    );
};

export default ThemeShow;
