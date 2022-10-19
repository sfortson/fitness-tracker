import styled from '@emotion/styled';
import { Button, Container, Grid, TextField, Typography } from '@mui/material';
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

const validationSchema = yup.object({
  email: yup.string().email('Enter a valid email').required('Email is required'),
  password: yup.string().min(8, 'Password should be of minimum 8 characters length').required('Password is required'),
});

export function Login() {
  const formik = useFormik({
    initialValues: {
      email: 'foobar@example.com',
      password: 'foobar',
    },
    validationSchema,
    onSubmit: (values) => {
      alert(JSON.stringify(values, null, 2));
    },
  });

  return (
    <Container>
      <Grid container>
        <Grid item xs={12}>
          <StyledTitle variant="h2">Login</StyledTitle>
        </Grid>
        <Grid item xs={12}>
          <StyledTitle variant="h3">Enter Username and Password</StyledTitle>
          <Grid item xs={4}>
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
