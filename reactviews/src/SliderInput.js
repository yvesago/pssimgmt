import React, {useState} from 'react';
import Slider from '@mui/material/Slider';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';

import { useInput } from 'react-admin';

const classes = {
    root: {
        width: 300,
    },
};


const SliderInput = (props) => {

    const {
        field: { name, value, onChange, ...rest },
        fieldState: { isTouched, invalid },
        formState: { isSubmitted }
    } = useInput(props);


    const [sliderValue, setSliderValue] = useState(value);

    const onSliderValueChange = (e, val) => {
        setSliderValue(val);
        //console.log(val);
        onChange(val);
    };

    const first = (a) => {
        return a[0].value;
    };

    const last = (a) => {
        return a[a.length-1].value;
    };

    return (
        <Box sx={classes.root}>
            <Typography id="discrete-slider-small-steps" gutterBottom={true}>
                {props.label}
            </Typography>
            <Slider
                name={name}
                value={Number(sliderValue)}
                aria-labelledby="discrete-slider-small-steps"
                step={null}
                valueLabelDisplay="auto"
                marks={props.choices}
                min={first(props.choices)}
                max={last(props.choices)}
                onChange={onSliderValueChange}
                error={(isTouched || isSubmitted) && invalid ? 'error' : 'false'}
                {...rest}
            />
        </Box>
    ); 
};

export default SliderInput;
