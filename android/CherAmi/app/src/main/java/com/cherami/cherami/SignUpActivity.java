package com.cherami.cherami;

import android.app.Activity;
import android.content.Intent;
import android.content.res.AssetManager;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.view.ViewPager;
import android.text.TextUtils;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.StringHttpMessageConverter;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestTemplate;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.UnsupportedEncodingException;
import java.net.MalformedURLException;
import java.net.URI;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;

import com.loopj.android.http.*;


public class SignUpActivity extends Activity {

    EditText mUsername;
    EditText mEmail;
    EditText mPassword;
    EditText mConfirmPassword;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_sign_up);
        getActionBar().hide();

        //Get handle, email, password, and confirm fields
        mUsername = (EditText) findViewById(R.id.username);
        mEmail = (EditText) findViewById(R.id.email);
        mPassword = (EditText) findViewById(R.id.password);
        mConfirmPassword = (EditText) findViewById(R.id.confirmPassword);

    }

    public void showLogin(View view) {
        Intent intent = new Intent(this, LoginActivity.class);
        startActivity(intent);
        finish();
    }

    public void attemptCreateAccount(View view) {
        AsyncHttpClient client = new AsyncHttpClient();

        JSONObject jsonParams = new JSONObject();
        try {
            jsonParams.put("handle", mUsername.getText().toString());
            jsonParams.put("email", mEmail.getText().toString());
            jsonParams.put("password", mPassword.getText().toString());
            jsonParams.put("confirmpassword", mConfirmPassword.getText().toString());
        } catch (JSONException j) {
            System.out.println("DONT LIKE JSON!");
        }
        StringEntity entity = null;
        try {
            entity = new StringEntity(jsonParams.toString());
        } catch (UnsupportedEncodingException i) {
            System.out.println("DONT LIKE TO STRING!");
        }

        client.post(this.getApplicationContext(), "http://10.0.1.29:8228/api/signup", entity, "application/json", new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                // called before request is started
                System.out.println("STARTING POST REQUEST");

            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                String s = new String(response);

                // called when response HTTP status is "200 OK"
                System.out.println("SUCCESS IN POSTING THAT USER!");
                System.out.println("status code: " + statusCode);
                System.out.println("response: " + s);
            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                String s = new String(errorResponse);
                // called when response HTTP status is "4XX" (eg. 401, 403, 404)
                System.out.println("AWE FUCK");
                System.out.println("status code: " + statusCode);
                System.out.println("response: " + s);
                e.printStackTrace();

            }

            @Override
            public void onRetry(int retryNo) {
                // called when request is retried
                System.out.println("RETRYING?!?!");
            }
        });



//        View focusView = null;
//        Boolean cancel = false;
//
//        String username = mUsername.getText().toString();
//        String email = mEmail.getText().toString();
//        String password = mPassword.getText().toString();
//        String confrimPassword = mConfirmPassword.getText().toString();
//        /* First, data sanitization: No fields should be left blank, email should have @ symbol,
//        password and Confirm password should be the same (this is done in back end)
//        Also, handle/username must be unique (also back end?)
//        Then,
//        POST to db with Handle, Email, Password, and Confirm
//        */
//
//        // Check for a valid email address.
//        if (TextUtils.isEmpty(email)) {
//            mEmail.setError(getString(R.string.error_field_required));
//            focusView = mEmail;
//            cancel = true;
//        } else if (!isEmailValid(email)) {
//            mEmail.setError(getString(R.string.error_invalid_email));
//            focusView = mEmail;
//            cancel = true;
//        }
//
//        if (TextUtils.isEmpty(password)) {
//            mPassword.setError(getString(R.string.error_field_required));
//            focusView = mPassword;
//            cancel = true;
//        } else if (!isPasswordValid(password)) {
//            mPassword.setError(getString(R.string.error_invalid_password));
//            focusView = mPassword;
//            cancel = true;
//        }
//
//        if (!confrimPassword.equals(password)) {
//            mConfirmPassword.setError(getString(R.string.error_invalid_confirm));
//            focusView = mConfirmPassword;
//            cancel = true;
//        }
//
//        if (TextUtils.isEmpty(username)) {
//            mUsername.setError(getString(R.string.error_field_required));
//            focusView = mUsername;
//            cancel = true;
//        }
//
//
//        if (cancel) {
//            //Something is wrong; don't sign up
//            focusView.requestFocus();
//        } else {
//            // Sign them up
//            new HttpRequestTask().execute(MediaType.APPLICATION_JSON);
//
//        }

    }


//    private class HttpRequestTask extends AsyncTask<MediaType, Void, String> {
//        private NewUser u;
//        private String host;
//
//        @Override
//        protected void onPreExecute() {
//            u = new NewUser(mUsername.getText().toString(), mEmail.getText().toString(), mPassword.getText().toString(), mConfirmPassword.getText().toString());
//
//            AssetManager assetManager = getResources().getAssets();
//            InputStream inputStream = null;
//            try {
//                inputStream = assetManager.open("config.properties");
//            } catch (IOException e) {
//                e.printStackTrace();
//            }
//            Properties properties = new Properties();
//            try {
//                properties.load(inputStream);
//            } catch (IOException e) {
//                e.printStackTrace();
//            }
//            host = properties.getProperty("myUrl");
//        }
//
//
//        @Override
//        protected String doInBackground(MediaType... params) {
//            try {
//                final String url = host + "/api/signup";
//
//
//                // Set the Content-Type header
//                HttpHeaders requestHeaders = new HttpHeaders();
//                requestHeaders.setContentType(MediaType.APPLICATION_JSON);
//                HttpEntity<NewUser> requestEntity = new HttpEntity<NewUser>(u, requestHeaders);
//
//                // Create a new RestTemplate instance
//                RestTemplate restTemplate = new RestTemplate();
//
//                // Add the Jackson and String message converters
//                restTemplate.getMessageConverters().add(new StringHttpMessageConverter());
//                restTemplate.getMessageConverters().add(new MappingJackson2HttpMessageConverter());
//
//
//                // Make the HTTP POST request, marshaling the request to JSON, and the response to a String
//                ResponseEntity<String> response = restTemplate.exchange(url, HttpMethod.POST, requestEntity, String.class);
//                return response.getBody();
//
//            } catch (Exception e) {
//                Log.e("MainActivity", e.getMessage(), e);
////                if (e instanceof MyCustomException) {
////                    MyCustomException exception = (MyCustomException) e;
////                    Log.d("Error: ", "An error occurred while calling api/user/account API endpoint: " + e.getMessage());
////                } else {
////                    Log.d("Error: ", "An error occurred while trying to parse Login Response JSON object");
////                }
//            }
//
//            return null;
//        }
//
//        @Override
//        protected void onPostExecute(String result) {
//            Toast toast = Toast.makeText(getApplicationContext(), result, Toast.LENGTH_LONG);
//            toast.show();
//        }
//    }

    private boolean isEmailValid(String email) {
        //TODO: Replace this with your own logic
        return email.contains("@");
    }

    private boolean isPasswordValid(String password) {
        //TODO: Replace this with your own logic
        return password.length() > 4;
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.sign_up, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();
        if (id == R.id.action_settings) {
            return true;
        }
        return super.onOptionsItemSelected(item);
    }
}
