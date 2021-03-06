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

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.UnsupportedEncodingException;

/**
 * Created by goalsman on 10/7/14.
 */
public class LoginActivity extends Activity {
    EditText mUsername;
    EditText mPassword;
    ProgressDialog dialog;
    Context context;
    TextView signUpButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        this.context = getApplicationContext();

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);
        getActionBar().hide();

        signUpButton = (TextView) findViewById(R.id.signUp);
        signUpButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                signUpButton.setTextColor(Color.parseColor("#4cc1ff"));
                showSignUp();
            }

        });

        mUsername = (EditText) findViewById(R.id.username);
        mPassword = (EditText) findViewById(R.id.password);
    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {

        return true;
    }

    public JSONObject getUserObjectRequestAsJson () {
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

    public void attemptLoginAccount() {
        AsyncHttpClient client = new AsyncHttpClient();

        client.post(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "sessions",
                convertJsonUserToStringEntity(getUserObjectRequestAsJson()), "application/json",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        dialog = ProgressDialog.show(LoginActivity.this, "",
                                "Loading. Please wait...", true);
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

                        Intent intent = new Intent(getApplicationContext(), MainActivity.class);
                        startActivity(intent);
                        finish();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                        dialog.dismiss();
                        String responseText = null;

                        try {
                            if (!ErrorHandle.isNetworkConnected(errorResponse)) {
                                ErrorHandle.displayNetworkErrorModal(LoginActivity.this);

                            } else {
                                responseText = new JSONObject(new String(errorResponse)).getString("reason");
                                Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                                toast.show();
                            }
                        } catch (JSONException j) {

                        }
                    }

                    @Override
                    public void onRetry(int retryNo) {

                    }
                });
    }

    public void showSignUp (){
        Intent intent = new Intent(this, SignUpActivity.class);
        startActivity(intent);
        finish();
    }

    public void loginClicked(View view) {
        View focusView = null;
        Boolean cancel = false;
        String username = mUsername.getText().toString();
        String password = mPassword.getText().toString();

        if (TextUtils.isEmpty(password)) {
            mPassword.setError(getString(R.string.error_field_required));
            focusView = mPassword;
            cancel = true;
        } else if (!isPasswordValid(password)) {
            mPassword.setError(getString(R.string.error_invalid_password));
            focusView = mPassword;
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
            // Sign them up; for now, redirect to Main
            attemptLoginAccount();
        }

    }

    private boolean isPasswordValid(String password) {
        return password.length() >= 8;
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
