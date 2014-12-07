package com.cherami.cherami;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.Toast;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.IOException;
import java.io.InputStream;
import java.io.UnsupportedEncodingException;
import java.util.Properties;

import com.loopj.android.http.*;


public class SignUpActivity extends Activity {

    EditText mUsername;
    EditText mEmail;
    EditText mPassword;
    EditText mConfirmPassword;
    SharedPreferences prefs;

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        Context context = getApplicationContext();
        System.out.println(context);

        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
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

    public JSONObject getUserObjectRequestAsJson () {
        JSONObject jsonParams = new JSONObject();
        try {
            jsonParams.put("handle", mUsername.getText().toString());
            jsonParams.put("email", mEmail.getText().toString());
            jsonParams.put("password", mPassword.getText().toString());
            jsonParams.put("confirmpassword", mConfirmPassword.getText().toString());
        } catch (JSONException j) {
            System.out.println(j);
        }
        System.out.println(jsonParams);
        return jsonParams;
    }
    public JSONObject getLoginObjectRequestAsJson () {
        JSONObject jsonParams = new JSONObject();
        try {
            jsonParams.put("handle", mUsername.getText().toString());
            jsonParams.put("password", mPassword.getText().toString());
        } catch (JSONException j) {
            System.out.println(j);
        }
        return jsonParams;
    }

    public StringEntity convertJsonUserToStringEntity (JSONObject jsonParams) {
        StringEntity entity = null;
        try {
            entity = new StringEntity(jsonParams.toString());
        } catch (UnsupportedEncodingException i) {
            System.out.println("DONT LIKE TO STRING!");
        }
        return entity;
    }

    public void getAuthToken() {
        AsyncHttpClient client = new AsyncHttpClient();

        client.post(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "sessions",
                convertJsonUserToStringEntity(getLoginObjectRequestAsJson()), "application/json",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        // called before request is started
                        System.out.println("STARTING POST REQUEST");

                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                        String s = new String(response);
                        System.out.println("HERES " + s);
                        JSONObject returnVal = new JSONObject();
                        try {
                            returnVal = new JSONObject(s);
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        try {
                            String sessionKey = "com.cherami.cherami.token";
                            String userKey = "com.cherami.cherami.username";

                            prefs.edit().putString(userKey, mUsername.getText().toString()).apply();
                            prefs.edit().putString(sessionKey, returnVal.getString("token")).apply();
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        // called when response HTTP status is "200 OK"

                        String responseText = null;
                        try {
                            responseText = new JSONObject(new String(response)).getString("response");
                        } catch (JSONException j) {
                            System.out.println(j);
                        }

                        Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();

                        Intent intent = new Intent(getApplicationContext(), MainActivity.class);
                        startActivity(intent);
                        finish();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                        // called when response HTTP status is "4XX" (eg. 401, 403, 404)
                        System.out.println("AWE RATS");

                        String responseText = null;
                        try {
                            responseText = new JSONObject(new String(errorResponse)).getString("reason");
                        } catch (JSONException j) {
                            System.out.println(j);
                        }

                        Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                        e.printStackTrace();
                    }

                    @Override
                    public void onRetry(int retryNo) {
                        // called when request is retried
                        System.out.println("RETRYING?!?!");
                    }
                });
    }
    public void attemptCreateAccount() {
        AsyncHttpClient client = new AsyncHttpClient();

        client.post(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "signup",
                    convertJsonUserToStringEntity(getUserObjectRequestAsJson()), "application/json",
                    new AsyncHttpResponseHandler() {

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

                String responseText = null;
                try {
                    responseText = new JSONObject(new String(response)).getString("response");
                } catch (JSONException j) {
                    System.out.println(j);
                }

                Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                toast.show();
                getAuthToken();
            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                // called when response HTTP status is "4XX" (eg. 401, 403, 404)

                String responseText = null;
                try {
                    if (errorResponse == null) {
                        new AlertDialog.Builder(SignUpActivity.this)
                                .setTitle("Network Error")
                                .setMessage("You're not connected to the network :(")
                                .setPositiveButton(getResources().getString(R.string.retry), new DialogInterface.OnClickListener() {
                                    public void onClick(DialogInterface dialog, int which) {
                                        // retry connection
                                        attemptCreateAccount();
                                    }
                                })
                                .setNegativeButton(android.R.string.ok, new DialogInterface.OnClickListener() {
                                    public void onClick(DialogInterface dialog, int which) {
                                        // do nothing
                                    }
                                })
                                .setIcon(android.R.drawable.ic_dialog_alert)
                                .show();
                    } else {
                        responseText = new JSONObject(new String(errorResponse)).getString("reason");
                        Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                    }

                } catch (JSONException j) {
                    System.out.println(j);
                }
                e.printStackTrace();
            }

            @Override
            public void onRetry(int retryNo) {
                // called when request is retried
                System.out.println("RETRYING?!?!");
            }
        });
    }

    public void signupButtonClicked (View view) {
        View focusView = null;
        Boolean cancel = false;

        String username = mUsername.getText().toString();
        String email = mEmail.getText().toString();
        String password = mPassword.getText().toString();
        String confirmPassword = mConfirmPassword.getText().toString();
        /* First, data sanitization: No fields should be left blank, email should have @ symbol,
        password and Confirm password should be the same (this is done in back end)
        Also, handle/username must be unique (also back end?)
        Then,
        POST to db with Handle, Email, Password, and Confirm
        */

        // Check for a valid email address.
        if (TextUtils.isEmpty(email)) {
            mEmail.setError(getString(R.string.error_field_required));
            focusView = mEmail;
            cancel = true;
        } else if (!isEmailValid(email)) {
            mEmail.setError(getString(R.string.error_invalid_email));
            focusView = mEmail;
            cancel = true;
        }

        if (TextUtils.isEmpty(password)) {
            mPassword.setError(getString(R.string.error_field_required));
            focusView = mPassword;
            cancel = true;
        } else if (!isPasswordValid(password)) {
            mPassword.setError(getString(R.string.error_invalid_password));
            focusView = mPassword;
            cancel = true;
        }

        if (!confirmPassword.equals(password)) {
            mConfirmPassword.setError(getString(R.string.error_invalid_confirm));
            focusView = mConfirmPassword;
            cancel = true;
        }

        if (TextUtils.isEmpty(username)) {
            mUsername.setError(getString(R.string.error_field_required));
            focusView = mUsername;
            cancel = true;
        }


        if (cancel) {
            //Something is wrong; don't sign up
            focusView.requestFocus();
        } else {
            // Attempt to sign them up
            attemptCreateAccount();
        }
    }

    private boolean isEmailValid(String email) {
        //TODO: Replace this with your own logic
        return email.contains("@");
    }

    private boolean isPasswordValid(String password) {
        //TODO: Replace this with your own logic
        return password.length() >= 8;
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {

        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();
        if (id == R.id.action_logout) {
            return true;
        }
        return super.onOptionsItemSelected(item);
    }
}
