package com.cherami.cherami;

import android.app.ProgressDialog;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.Button;
import android.widget.ListView;
import android.widget.Spinner;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;


/**
 * A simple {@link Fragment} subclass.
 * Activities that contain this fragment must implement the
 * {@link Circles.OnFragmentInteractionListener} interface
 * to handle interaction events.
 * Use the {@link Circles#newInstance} factory method to
 * create an instance of this fragment.
 *
 */
public class Circles extends Fragment {
    // TODO: Rename parameter arguments, choose names that match
    // the fragment initialization parameters, e.g. ARG_ITEM_NUMBER
    private static final String ARG_PARAM1 = "param1";
    private static final String ARG_PARAM2 = "param2";
    private ListView circleList;
    // TODO: Rename and change types of parameters
    private String mParam1;
    private String mParam2;
    private Spinner spinner;
    SharedPreferences prefs;
    CircleAdapter adapter;
    Button createNewCircleButton;
    Button refreshButton;
    Context context;
    ProgressDialog dialog;

    public Circles() {

    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        this.context = getActivity().getApplicationContext();
        super.onCreate(savedInstanceState);
    }

    public void filterCircles () {
        String value = spinner.getSelectedItem().toString();
        final CircleAdapter newCircleAdapter;

        try {
            Circle [] circleData = Circles.this.adapter.getData();

            if (value.equals("All")) {
                newCircleAdapter = new CircleAdapter(getActivity(),
                        R.layout.circle_item_row, circleData);
            } else {
                ArrayList<Circle> filteredCircles = new ArrayList<Circle>();

                for (int i = 0; i < circleData.length; i++) {
                    if (circleData[i].getVisibility().equals(value.toLowerCase())) {
                        filteredCircles.add(circleData[i]);
                    }
                }

                Circle [] newCircleData = filteredCircles.toArray(new Circle[filteredCircles.size()]);
                newCircleAdapter = new CircleAdapter(getActivity(),
                        R.layout.circle_item_row, newCircleData);
            }

            circleList = (ListView) Circles.this.getView().findViewById(R.id.circleList);
            circleList.setAdapter(newCircleAdapter);
            circleList.setOnItemClickListener(new AdapterView.OnItemClickListener() {
                @Override
                public void onItemClick(AdapterView<?> parent, View view, int position,
                                        long id) {
                    Intent intent = new Intent(getActivity().getApplicationContext(), CircleResult.class);
                    Bundle mBundle = new Bundle();
                    try {
                        mBundle.putString("owner",newCircleAdapter.getItem(position).getCircle().getString("owner"));
                        mBundle.putString("circleName", newCircleAdapter.getItem(position).getCircle().getString("name"));
                        mBundle.putString("joinVisibility", "none");
                        mBundle.putString("circleid", newCircleAdapter.getItem(position).getCircle().getString("url"));
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                    intent.putExtras(mBundle);
                    startActivity(intent);
                }
            });

        } catch (NullPointerException n) {

        }

    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {

        View rootView = inflater.inflate(R.layout.fragment_circles, container, false);
        getCircles(rootView);

        spinner = (Spinner) rootView.findViewById(R.id.filter_spinner);
        createNewCircleButton = (Button) rootView.findViewById(R.id.createNewCircle);
        refreshButton = (Button) rootView.findViewById(R.id.refreshButton);

        spinner.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                filterCircles();
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {

            }
        });

        createNewCircleButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                displayCreateCircleModal();
            }
        });

        refreshButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Circles.this.getCircles(Circles.this.getView());
            }
        });

        return rootView;
    }

    public void displayCreateCircleModal () {
        CreateCircleModal createCircleModalFragment = new CreateCircleModal();
        createCircleModalFragment.show(getFragmentManager(), "dialog");
    }

    public void setCircleAdapter (CircleAdapter circleAdapter) {
        this.adapter = circleAdapter;
    }

    public void getCircles(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);
        String username = ApiHelper.getUsername(context);
        RequestParams params = new RequestParams();
        params.put("user", username);

        client.addHeader("Authorization", token);
        client.get(context, ApiHelper.getLocalUrlForApi(getResources()) + "circles",
                   params, new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                dialog = ProgressDialog.show(getActivity(), "",
                        "Loading. Please wait...", true);
            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;

                try {
                    responseText = new JSONObject(new String(responseBody)).getString("results");
                    JSONArray y = new JSONArray(responseText);
                    Circle circle_data[] = new Circle[y.length()];

                    for (int x = 0; x < y.length(); x++) {

                        circle_data[x] = new Circle(new JSONObject(y.get(x).toString()));
                    }

                    final CircleAdapter adapter = new CircleAdapter(getActivity(),
                            R.layout.circle_item_row, circle_data);
                    Circles.this.setCircleAdapter(adapter);
                    circleList = (ListView) view.findViewById(R.id.circleList);
                    circleList.setAdapter(adapter);

                } catch (JSONException j) {

                }

                filterCircles();
                dialog.dismiss();
            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                dialog.dismiss();
                String responseText = null;

                try {
                    if (!ErrorHandle.isNetworkConnected(errorResponse)) {
                        ErrorHandle.displayNetworkErrorModal(getActivity());

                    } else {
                        responseText = new JSONObject(new String(errorResponse)).getString("reason");
                        if (ErrorHandle.isTokenExpired(responseText)) {
                            ErrorHandle.displayTokenModal(getActivity());
                        }
                    }
                } catch (JSONException j) {

                }

            }
        });
    }

}
