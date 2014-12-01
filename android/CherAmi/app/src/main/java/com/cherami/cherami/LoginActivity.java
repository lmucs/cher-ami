package com.cherami.cherami;

import android.app.ActionBar;
import android.app.Activity;
import android.app.AlertDialog;
import android.app.Fragment;
import android.app.FragmentManager;
import android.app.FragmentTransaction;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.preference.PreferenceManager;
import android.support.v13.app.FragmentPagerAdapter;
import android.support.v4.view.ViewPager;
import android.text.TextUtils;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.EditText;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.io.UnsupportedEncodingException;
import java.util.Locale;
import java.util.Properties;

/**
 * Created by goalsman on 10/7/14.
 */
public class LoginActivity extends Activity {
    /**
     * The {@link android.support.v4.view.PagerAdapter} that will provide
     * fragments for each of the sections. We use a
     * {@link android.support.v13.app.FragmentPagerAdapter} derivative, which will keep every
     * loaded fragment in memory. If this becomes too memory intensive, it
     * may be best to switch to a
     * {@link android.support.v13.app.FragmentStatePagerAdapter}.
     */

    /**
     * The {@link android.support.v4.view.ViewPager} that will host the section contents.
     */
    EditText mUsername;
    EditText mPassword;

    SharedPreferences prefs;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Context context = getApplicationContext();
        System.out.println(context);

        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);
        getActionBar().hide();

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
            System.out.println("DONT LIKE JSON!");
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

    public void attemptLoginAccount() {
        AsyncHttpClient client = new AsyncHttpClient();

        client.post(this.getApplicationContext(), "http://" + ApiHelper.getLocalUrlForApi(getResources()) + "/api/sessions",
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
                        JSONObject returnVal = new JSONObject();
                        try {
                            returnVal = new JSONObject(s);
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        try {
                            String sessionKey = "com.cherami.cherami.token";
                            String userKey = "com.cherami.cherami.username";

                            prefs.edit().putString(sessionKey, returnVal.getString("token")).apply();
                            prefs.edit().putString(userKey, mUsername.getText().toString()).apply();
                            System.out.println("Token " + prefs.getString(sessionKey, null));
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        // called when response HTTP status is "200 OK"

                        String responseText = null;
                        try {

                            responseText = new JSONObject(new String(response)).getString("response");
                        } catch (JSONException j) {
                            System.out.println("Dont like JSON");
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

                        String responseText = null;
                        try {
                            if (errorResponse == null) {
                                new AlertDialog.Builder(LoginActivity.this)
                                        .setTitle("Network Error")
                                        .setMessage("You're not connected to the network :(")
                                        .setPositiveButton(getResources().getString(R.string.retry), new DialogInterface.OnClickListener() {
                                            public void onClick(DialogInterface dialog, int which) {
                                                // retry connection
                                                attemptLoginAccount();
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
                            System.out.println("Dont like JSON");
                        }
                    }

                    @Override
                    public void onRetry(int retryNo) {
                        // called when request is retried
                        System.out.println("RETRYING?!?!");
                    }
                });
    }

    public void goBackToSignUp (View view){
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
        //TODO: Replace this with your own logic
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
