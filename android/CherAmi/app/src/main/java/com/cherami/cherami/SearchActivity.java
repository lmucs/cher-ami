package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ListView;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;


public class SearchActivity extends Activity {
    private ListView userList;
    SharedPreferences prefs;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Context context = this.getApplicationContext();
        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_search);
        getActionBar().hide();
    }

    public void getUsers(View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String sessionKey = "com.cherami.cherami.token";
        String token = prefs.getString(sessionKey, null);
        String searchInput = ((EditText)findViewById(R.id.search_bar)).getText().toString();
        System.out.println("Searching: " + searchInput);
        RequestParams params = new RequestParams();
        params.put("nameprefix", searchInput);
        params.put("sort", "joined");

        client.addHeader("Authorization", token);
        client.get(this.getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "users", params, new AsyncHttpResponseHandler() {

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
                    System.out.println(responseText);
                    JSONArray y = new JSONArray(responseText);
                    User user_data[] = new User[y.length()];
                    for (int x = 0; x < y.length(); x++) {
                        user_data[x] = new User(new JSONObject(y.get(x).toString()));
                    }

                    final UserAdapter adapter = new UserAdapter(SearchActivity.this,
                            R.layout.user_item_row, user_data);

                    System.out.println("adatper " + adapter);


                    userList = (ListView) findViewById(R.id.searchList);

                    System.out.println("userlist: " + userList);

                    userList.setAdapter(adapter);
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
        int id = item.getItemId();
        if (id == R.id.action_settings) {
            return true;
        }
        return super.onOptionsItemSelected(item);
    }
}
