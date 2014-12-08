package com.cherami.cherami;

import android.app.ActionBar;
import android.app.Activity;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;


public class SearchActivity extends Activity {
    private ListView userList;
    Context context;
    ProgressDialog dialog;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        this.context = this.getApplicationContext();
        super.onCreate(savedInstanceState);

        ActionBar actionBar = getActionBar();
        actionBar.setDisplayHomeAsUpEnabled(true);

        // Set view to XML file
        setContentView(R.layout.activity_search);
    }

    public void getUsers(View view) {
        AsyncHttpClient client = new AsyncHttpClient();

        //Static call to get token
        String token = ApiHelper.getSessionToken(context);
        String searchInput = ((EditText)findViewById(R.id.search_bar)).getText().toString();
        RequestParams params = new RequestParams();
        // Add search to HTTP call, sort by joined date in descending order
        params.put("nameprefix", searchInput);
        params.put("sort", "joined");

        client.addHeader("Authorization", token);
        client.get(context, ApiHelper.getLocalUrlForApi(getResources()) + "users", params, new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                // Display spinner
                dialog = ProgressDialog.show(SearchActivity.this, "",
                        "Loading. Please wait...", true);
            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                dialog.dismiss();
                String responseText = null;

                try {
                    responseText = new JSONObject(new String(responseBody)).getString("results");
                    JSONArray y = new JSONArray(responseText);
                    User user_data[] = new User[y.length()];

                    for (int x = 0; x < y.length(); x++) {
                        user_data[x] = new User(new JSONObject(y.get(x).toString()));
                    }

                    final UserAdapter adapter = new UserAdapter(SearchActivity.this,
                            R.layout.user_item_row, user_data);
                    userList = (ListView) findViewById(R.id.searchList);
                    userList.setAdapter(adapter);

                    // Prepare for click on username to nav to that user's profile
                    userList.setOnItemClickListener(new AdapterView.OnItemClickListener() {
                        @Override
                        public void onItemClick(AdapterView<?> parent, View view, int position,
                                                long id) {
                            Intent intent = new Intent(SearchActivity.this, OtherUserProfileActivity.class);
                            Bundle mBundle = new Bundle();
                            try {
                                mBundle.putString("handle",adapter.getItem(position).getUserName().getString("u.handle"));
                            } catch (JSONException e) {
                                e.printStackTrace();
                            }
                            intent.putExtras(mBundle);
                            startActivity(intent);
                        }
                    });
                } catch (JSONException j) {

                }

            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                dialog.dismiss();
                String responseText = null;

                try {
                    if (!NetworkCheck.isConnected(errorResponse)) {
                        NetworkCheck.displayNetworkErrorModal(SearchActivity.this);

                    } else {
                        responseText = new JSONObject(new String(errorResponse)).getString("reason");
                        Toast toast = Toast.makeText(getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                    }
                } catch (JSONException j) {

                }

            }
        });
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
}
