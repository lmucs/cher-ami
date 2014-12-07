package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.RadioButton;
import android.widget.TextView;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.UnsupportedEncodingException;


public class CircleResult extends Activity {

    SharedPreferences prefs;
    TextView textElement;
    String circleName;
    String owner;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Context context = this.getApplicationContext();
        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_circle_result);

        textElement=(TextView)findViewById(R.id.circleName);
        Bundle recdData = getIntent().getExtras();
        circleName = recdData.getString("circleName");
        owner = recdData.getString("owner");
        View joinButton = findViewById(R.id.joinCircle);
        if(recdData.getString("joinVisibility").equals("none")){
            joinButton.setVisibility(View.GONE);
        }
        textElement.setText(circleName);
    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.circle_result, menu);
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

    public JSONObject getJoinParamsAsJson () {
        JSONObject jsonParams = new JSONObject();

        try {
            jsonParams.put("Circle", circleName);
            jsonParams.put("Target", owner);
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

    public void joinCircle(View view){

            AsyncHttpClient client = new AsyncHttpClient();
            String sessionKey = "com.cherami.cherami.token";
            String token = prefs.getString(sessionKey, null);

            client.addHeader("Authorization", token);
            client.post(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "join",
                    convertJsonUserToStringEntity(getJoinParamsAsJson()), "application/json",
                    new AsyncHttpResponseHandler() {

                        @Override
                        public void onStart() {

                        }

                        @Override
                        public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                            String responseText = null;

                            try {
                                responseText = new JSONObject(new String(response)).getString("response");
                            } catch (JSONException j) {

                            }
                            
                            Toast toast = Toast.makeText(CircleResult.this.getApplicationContext(), responseText, Toast.LENGTH_LONG);
                            toast.show();

                        }

                        @Override
                        public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                            String responseText = null;

                            try {
                                responseText = new JSONObject(new String(errorResponse)).getString("Reason");
                            } catch (JSONException j) {

                            }

                            Toast toast = Toast.makeText(CircleResult.this.getApplicationContext(), responseText, Toast.LENGTH_LONG);
                            toast.show();
                            e.printStackTrace();
                        }

                        @Override
                        public void onRetry(int retryNo) {

                        }
                    });
    }
}
