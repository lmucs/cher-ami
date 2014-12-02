package com.cherami.cherami;

import android.app.ActionBar;
import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.ListView;
import android.widget.TextView;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

/**
 * Created by Geoff on 11/29/2014.
 */
public class OtherUserProfileActivity extends Activity{

    private ListView circleList;
    SharedPreferences prefs;
    TextView textElement;
    String myVal;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Context context = this.getApplicationContext();
        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
        ActionBar actionBar = getActionBar();
        actionBar.setDisplayHomeAsUpEnabled(true);
        setContentView(R.layout.other_user_profile);

        textElement=(TextView)findViewById(R.id.otherUsername);
        Bundle recdData = getIntent().getExtras();
        myVal = recdData.getString("handle");
        textElement.setText(myVal);
        getOtherUserCircles(this.findViewById(R.id.otherCircleFeed).getRootView());
        System.out.println(myVal);
    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.search, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        switch (item.getItemId())
        {
            case android.R.id.home:
                Intent intent = new Intent(this, MainActivity.class);
                intent.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP | Intent.FLAG_ACTIVITY_SINGLE_TOP);
                startActivity(intent);
                return true;
            default:
                return super.onOptionsItemSelected(item);
        }
    }

    public String processDate(String date){
        return date.substring(0, date.lastIndexOf("T"));
    }

    public String getLocalUrlForApi () {
        AssetManager assetManager = getResources().getAssets();
        InputStream inputStream = null;
        try {
            inputStream = assetManager.open("config.properties");
        } catch (IOException e) {
            e.printStackTrace();
        }
        Properties properties = new Properties();
        try {
            properties.load(inputStream);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return properties.getProperty("myUrl");
    }

    public void getOtherUserCircles(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String sessionKey = "com.cherami.cherami.token";
        String token = prefs.getString(sessionKey, null);

        client.addHeader("Authorization", token);
        client.get(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "circles?user=" + myVal, new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                // called before request is started
                System.out.println("STARTING GET REQUEST");

            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;
                try {
                    responseText = new JSONObject(new String(responseBody)).getString("results");
                    JSONArray y = new JSONArray(responseText);
                    OtherUserCircle circle_data[] = new OtherUserCircle[y.length()];
                    for (int x = 0; x < y.length(); x++) {

                        circle_data[x] = new OtherUserCircle(new JSONObject(y.get(x).toString()).getString("name"), new JSONObject(y.get(x).toString()).getString("owner"), processDate(new JSONObject(y.get(x).toString()).getString("created")));                    }

                    OtherUserProfileAdapter adapter = new OtherUserProfileAdapter(OtherUserProfileActivity.this,
                            R.layout.other_user_circle_row, circle_data);


                    circleList = (ListView) view.findViewById(R.id.otherCircleFeed);

                    circleList.setAdapter(adapter);
                } catch (JSONException j) {
                    System.out.println(j);
                }

            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                // called when response HTTP status is "4XX" (eg. 401, 403, 404)

                String responseText = null;
                try {
                    responseText = new JSONObject(new String(errorResponse)).getString("reason");

                } catch (JSONException j) {
                    System.out.println(j);
                }

            }
        });
    }
}
