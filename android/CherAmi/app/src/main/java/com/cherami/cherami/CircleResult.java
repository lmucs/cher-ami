package com.cherami.cherami;

import android.app.ActionBar;
import android.app.Activity;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.RadioButton;
import android.widget.TextView;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.UnsupportedEncodingException;


public class CircleResult extends Activity {
    private ListView feedList;
    TextView textElement;
    String circleName;
    String owner;
    ProgressDialog dialog;
    Context context;
    FeedAdapter adapter;
    Bundle recdData;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        this.context = this.getApplicationContext();
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_circle_result);
        ActionBar actionBar = getActionBar();
        actionBar.setDisplayHomeAsUpEnabled(true);

        textElement=(TextView)findViewById(R.id.circleName);
        this.recdData = getIntent().getExtras();
        circleName = recdData.getString("circleName");

        owner = recdData.getString("owner");
        View joinButton = findViewById(R.id.joinCircle);
        if(recdData.getString("joinVisibility").equals("none")){
            joinButton.setVisibility(View.GONE);
        }
        textElement.setText(circleName);
        getFeed(this.findViewById(R.id.feedList).getRootView());
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
        switch (item.getItemId()) {
            case android.R.id.home:
                Intent intent = new Intent(this, MainActivity.class);
                intent.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP | Intent.FLAG_ACTIVITY_SINGLE_TOP);
                startActivity(intent);
                return true;
            default:
                return super.onOptionsItemSelected(item);
        }
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

    public void setFeedAdapter (FeedAdapter adapter) {
        this.adapter = adapter;
    }

    public void getFeed(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);
        RequestParams params = new RequestParams();
        String circleURL = this.recdData.getString("circleid");
        params.put("circleid", circleURL.substring(circleURL.lastIndexOf("/") + 1));

        client.addHeader("Authorization", token);
        client.get(context, ApiHelper.getLocalUrlForApi(getResources()) + "messages", params,
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        dialog = ProgressDialog.show(CircleResult.this, "",
                                "Loading. Please wait...", true);
                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                        JSONArray responseText;

                        try {
                            responseText = new JSONArray(new String(responseBody));
                            FeedItem feed_data[] = new FeedItem[responseText.length()];

                            for (int x = 0; x < responseText.length(); x++) {
                                feed_data[x] = new FeedItem(new JSONObject(responseText.get(x).toString()));
                            }

                            final FeedAdapter adapter = new FeedAdapter(CircleResult.this,
                                    R.layout.feed_item_row, feed_data);
                            CircleResult.this.setFeedAdapter(adapter);
                            feedList = (ListView) view.findViewById(R.id.feedList);
                            feedList.setAdapter(adapter);

                        } catch (JSONException j) {
                            System.out.println(j);
                        }
                        dialog.dismiss();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                        dialog.dismiss();
                        String responseText = null;

                        try {
                            if (!NetworkCheck.isConnected(errorResponse)) {
                                NetworkCheck.displayNetworkErrorModal(CircleResult.this);

                            } else {
                                responseText = new JSONObject(new String(errorResponse)).getString("reason");
                                Toast toast = Toast.makeText(CircleResult.this.getApplicationContext(), responseText, Toast.LENGTH_LONG);
                                toast.show();
                            }
                        } catch (JSONException j) {

                        }

                    }
                });
    }

    public void joinCircle(View view){
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(this.context);

        client.addHeader("Authorization", token);
        client.post(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "join",
                convertJsonUserToStringEntity(getJoinParamsAsJson()), "application/json",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        dialog = ProgressDialog.show(CircleResult.this, "",
                                "Loading. Please wait...", true);
                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                        dialog.dismiss();
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
                        dialog.dismiss();
                        String responseText = null;

                        try {
                            if (!NetworkCheck.isConnected(errorResponse)) {
                                NetworkCheck.displayNetworkErrorModal(CircleResult.this);

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
}
