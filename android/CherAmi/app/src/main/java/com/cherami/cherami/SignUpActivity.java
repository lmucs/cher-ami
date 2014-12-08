package com.cherami.cherami;

import android.app.Activity;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.UnsupportedEncodingException;

import com.loopj.android.http.*;


public class SignUpActivity extends Activity {

    EditText mUsername;
    EditText mEmail;
    EditText mPassword;
    EditText mConfirmPassword;
    ProgressDialog dialog;
    Context context;
    TextView loginButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        this.context = getApplicationContext();
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_sign_up);
        getActionBar().hide();

        loginButton = (TextView) findViewById(R.id.login);
        loginButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                loginButton.setTextColor(Color.parseColor("#4cc1ff"));
                showLogin();
            }

        });



        //Get handle, email, password, and confirm fields
        mUsername = (EditText) findViewById(R.id.username);
        mEmail = (EditText) findViewById(R.id.email);
        mPassword = (EditText) findViewById(R.id.password);
        mConfirmPassword = (EditText) findViewById(R.id.confirmPassword);

    }

    public void showLogin() {
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

        }
        return jsonParams;
    }
    public JSONObject getLoginObjectRequestAsJson () {
        JSONObject jsonParams = new JSONObject();

        try {
            jsonParams.put("handle", mUsername.getText().toString());
            jsonParams.put("password", mPassword.getText().toString());
        } catch (JSONException j) {

        }

        return jsonParams;
    }

    public StringEntity convertJsonUserToStringEntity (JSONObject jsonParams) {
        StringEntity entity = null;

        try {
            entity = new StringEntity(jsonParams.toString());
        } catch (UnsupportedEncodingException i) {

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

                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                        dialog.dismiss();
                        String s = new String(response);
                        JSONObject returnVal = new JSONObject();

                        try {
                            returnVal = new JSONObject(s);
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }

                        try {
                            ApiHelper.saveAuthorizationToken(context, mUsername.getText().toString(),
                                                             returnVal.getString("token"));
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }

                        String responseText = null;

                        try {
                            responseText = new JSONObject(new String(response)).getString("response");
                        } catch (JSONException j) {

                        }

                        Toast toast = Toast.makeText(SignUpActivity.this.context, responseText, Toast.LENGTH_LONG);
                        toast.show();

                        Intent intent = new Intent(SignUpActivity.this.context, MainActivity.class);
                        startActivity(intent);
                        finish();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                        dialog.dismiss();

                        String responseText = null;

                        try {
                            if (!NetworkCheck.isConnected(errorResponse)) {
                                NetworkCheck.displayNetworkErrorModal(SignUpActivity.this);

                            } else {
                                responseText = new JSONObject(new String(errorResponse)).getString("reason");
                                Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                                toast.show();
                            }
                        } catch (JSONException j) {

                        }

                        Toast toast = Toast.makeText(SignUpActivity.this.context, responseText, Toast.LENGTH_LONG);
                        toast.show();
                        e.printStackTrace();
                    }

                    @Override
                    public void onRetry(int retryNo) {

                    }
                });
    }
    public void attemptCreateAccount() {
        AsyncHttpClient client = new AsyncHttpClient();

        client.post(context, ApiHelper.getLocalUrlForApi(getResources()) + "signup",
                    convertJsonUserToStringEntity(getUserObjectRequestAsJson()), "application/json",
                    new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                dialog = ProgressDialog.show(SignUpActivity.this, "",
                        "Loading. Please wait...", true);
            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                String s = new String(response);
                String responseText = null;

                try {
                    responseText = new JSONObject(new String(response)).getString("response");
                } catch (JSONException j) {

                }

                Toast toast = Toast.makeText(SignUpActivity.this.context, responseText, Toast.LENGTH_LONG);
                toast.show();
                getAuthToken();
            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                dialog.dismiss();
                String responseText = null;

                try {
                    if (!NetworkCheck.isConnected(errorResponse)) {
                        NetworkCheck.displayNetworkErrorModal(SignUpActivity.this);

                    } else {
                        responseText = new JSONObject(new String(errorResponse)).getString("reason");
                        Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                    }
                } catch (JSONException j) {

                }
                e.printStackTrace();
            }

            @Override
            public void onRetry(int retryNo) {

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
            mEmail.setError("This email is not properly formatted");
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
        } else if (!isUsernameValid(username)) {
            mUsername.setError(getString(R.string.error_invalid_email));
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

    private boolean isUsernameValid (String username) {
        return username.matches("^[\\p{L}\\p{M}][\\d\\p{L}\\p{M}]*$");
    }

    private boolean isEmailValid(String email) {

        return email.matches("^\\w+([-+.']\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$");
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
