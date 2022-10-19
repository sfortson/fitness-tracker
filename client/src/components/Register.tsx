import styled from '@emotion/styled';
import {
  Button,
  Container,
  FormControl,
  FormControlLabel,
  FormLabel,
  Grid,
  Radio,
  RadioGroup,
  TextField,
  Typography,
} from '@mui/material';
import { useFormik } from 'formik';
import { rem } from 'polished';
import React from 'react';
import * as yup from 'yup';

const StyledTitle = styled(Typography)`
  padding-bottom: ${rem(36)};
`;

const StyledTextField = styled(TextField)`
  min-width: 100%;
  padding-bottom: ${rem(8)};
`;

const StyledFormControl = styled(FormControl)`
  padding-bottom: ${rem(8)};
`;

const validationSchema = yup.object({
  email: yup.string().email('Enter a valid email').required('Email is required'),
  password: yup.string().min(8, 'Password should be of minimum 8 characters length').required('Password is required'),
});

export function RegisterUser() {
  const formik = useFormik({
    initialValues: {
      email: 'foobar@example.com',
      password: 'foobar',
      birthDate: '2022-10-05',
      gender: 'female',
    },
    validationSchema,
    onSubmit: (values) => {
      console.log(values);
    },
  });

  return (
    <Container>
      <Grid container>
        <Grid item xs={12}>
          <StyledTitle variant="h2">Register New User</StyledTitle>
        </Grid>
        <Grid item xs={12}>
          <StyledTitle variant="h3">Enter Information</StyledTitle>
          <Grid item xs={5}>
            <div>
              <form onSubmit={formik.handleSubmit}>
                <StyledTextField
                  fullWidth
                  id="email"
                  name="email"
                  label="Email"
                  value={formik.values.email}
                  onChange={formik.handleChange}
                  error={formik.touched.email && Boolean(formik.errors.email)}
                  helperText={formik.touched.email && formik.errors.email}
                />
                <StyledFormControl>
                  <FormLabel id="demo-radio-buttons-group-label">Gender</FormLabel>
                  <RadioGroup
                    row
                    aria-labelledby="demo-radio-buttons-group-label"
                    defaultValue="female"
                    name="gender"
                    onChange={(evt) => formik.handleChange(evt)}
                  >
                    <FormControlLabel value="female" control={<Radio />} label="Female" />
                    <FormControlLabel value="male" control={<Radio />} label="Male" />
                  </RadioGroup>
                </StyledFormControl>
                <StyledFormControl>
                  <FormLabel id="birthdate-label">Birthdate</FormLabel>
                  <TextField
                    value={formik.values.birthDate}
                    name="birthDate"
                    id="birthDate"
                    type="date"
                    onChange={formik.handleChange}
                  />
                </StyledFormControl>
                <StyledTextField
                  fullWidth
                  id="password"
                  name="password"
                  label="Password"
                  type="password"
                  value={formik.values.password}
                  onChange={formik.handleChange}
                  error={formik.touched.password && Boolean(formik.errors.password)}
                  helperText={formik.touched.password && formik.errors.password}
                />
                <Button color="primary" variant="contained" type="submit">
                  Submit
                </Button>
              </form>
            </div>
          </Grid>
        </Grid>
      </Grid>
    </Container>
  );
}
