import React from 'react';

import DashboardCmp from './DashboardCmp';
import MenuDomaine from './MenuDomaine';
import { Card, Grid, Box } from '@mui/material';

import { RadarChart } from './RadarChart';


const HomeShow = () => {
    return (
        <Box sx={{ flexGrow: 1 }} >
            <Grid container spacing={2}>
                <Grid item xs={3}>
                    <Card><MenuDomaine /></Card>
                </Grid>
                <Grid item xs={9}>
                    <>
                        <div style={{right: '0', width:'400px', 'zIndex': -3, position: 'absolute'}}>
                            <RadarChart  />
                        </div>
                        <DashboardCmp />
                    </>
                </Grid>
            </Grid>
        </Box>
    );
};

export default HomeShow;

