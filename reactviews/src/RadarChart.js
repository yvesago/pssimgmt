import React from 'react';
import { Radar } from 'react-chartjs-2';
import { useStore } from 'react-admin';

import { Chart as ChartJS, RadialLinearScale, PointElement, LineElement, Filler, Title, Legend } from 'chart.js';
ChartJS.register( RadialLinearScale, PointElement, LineElement, Filler, Title, Legend);

const RadarData = {
    labels: [],
    datasets: [
        {
            label: 'actuelle',
            backgroundColor: 'rgba(202, 202, 236, .1)',
            borderColor: 'rgba(234, 136, 1, 1)',
            pointBackgroundColor: 'rgba(234, 136, 1, 1)',
            poingBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgba(234, 136, 1, 1)',
            // data: [7, 10, 8, 6, 5, 8]
        },
        {
            label: 'attendue',
            backgroundColor: 'rgba(202, 202, 236, .1)',
            borderColor: 'rgba(34, 202, 136, 1)',
            pointBackgroundColor: 'rgba(34, 202, 136, 1)',
            poingBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgba(34, 202, 136, 1)',
            // data: [8, 10, 9, 7, 6, 9]
        },
    ]
};

const RadarOptions = {
    plugins: {
        title: {
            display: true,
            text: 'Conformité',
            padding: {
                top: 50,
                bottom: 0
            },
        },
        legend: {
            position: 'bottom',
            align: 'start',
        },
    },
    scales: {
        r : {
            min: 0,
            max: 10,
            ticks: {
                stepSize: 2,
                showLabelBackdrop: false,
                backdropPadding: 8,
            },
            grid: {
                color: 'lightgrey',
            },
            angleLines: {
                color: 'lightgrey',
            },
        }
    },
};


export const RadarChart = () => {
    const [dom] = useStore('dom', { select: {id:0 }, evaluation: {}, mask: {etu:false, nok:false} });

    if (dom && dom.evaluation && dom.evaluation.axes && dom.select.id !== '0' && dom.select.name !== undefined) {
        RadarData.labels = dom.evaluation.axes;
        RadarData.datasets[0].data = dom.evaluation.conforme;
        RadarData.datasets[1].data = dom.evaluation.evolution;
        RadarOptions.plugins.title.text = 'Conformité: ' + dom.select.name;
        return (
            <Radar data={RadarData} options={RadarOptions} />
        );
    } else {
        return '';
    }
};
