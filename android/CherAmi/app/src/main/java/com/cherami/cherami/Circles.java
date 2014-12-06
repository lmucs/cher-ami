package com.cherami.cherami;

import android.app.Activity;
import android.app.FragmentManager;
import android.content.Context;
import android.content.DialogInterface;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.net.Uri;
import android.os.Bundle;
import android.app.Fragment;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListAdapter;
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

import java.io.IOException;
import java.io.InputStream;
import java.lang.reflect.Array;
import java.util.ArrayList;
import java.util.Properties;


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

    private OnFragmentInteractionListener mListener;

    /**
     * Use this factory method to create a new instance of
     * this fragment using the provided parameters.
     *
     * @param param1 Parameter 1.
     * @param param2 Parameter 2.
     * @return A new instance of fragment Circles.
     */
    // TODO: Rename and change types and number of parameters
    public static Circles newInstance(String param1, String param2) {
        Circles fragment = new Circles();
        Bundle args = new Bundle();
        args.putString(ARG_PARAM1, param1);
        args.putString(ARG_PARAM2, param2);
        fragment.setArguments(args);
        return fragment;
    }

    public Circles() {
        // Required empty public constructor
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        Context context = getActivity().getApplicationContext();
        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
        if (getArguments() != null) {
            mParam1 = getArguments().getString(ARG_PARAM1);
            mParam2 = getArguments().getString(ARG_PARAM2);
        }
    }

    public void filterCircles () {
        String value = spinner.getSelectedItem().toString();
        CircleAdapter newCircleAdapter;

        try {
            Circle [] circleData = Circles.this.adapter.getData();

            if (value.equals("All")) {
                newCircleAdapter = new CircleAdapter(getActivity(),
                        R.layout.circle_item_row, circleData);
            } else {
                ArrayList<Circle> filteredCircles = new ArrayList<Circle>();
                for (int i = 0; i < circleData.length; i++) {
                    if (circleData[i].visibility.equals(value.toLowerCase())) {
                        filteredCircles.add(circleData[i]);
                    }
                }

                Circle [] newCircleData = filteredCircles.toArray(new Circle[filteredCircles.size()]);
                newCircleAdapter = new CircleAdapter(getActivity(),
                        R.layout.circle_item_row, newCircleData);
            }

            circleList = (ListView) Circles.this.getView().findViewById(R.id.circleList);
            circleList.setAdapter(newCircleAdapter);

        } catch (NullPointerException n) {
            System.out.println("NO DATA YET!");
        }

    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {

        View rootView = inflater.inflate(R.layout.fragment_circles, container, false);
        getCircles(rootView);

        spinner = (Spinner) rootView.findViewById(R.id.filter_spinner);
        createNewCircleButton = (Button) rootView.findViewById(R.id.createNewCircle);

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

        return rootView;
    }

    public void displayCreateCircleModal () {

        CreateCircleModal createCircleModalFragment = new CreateCircleModal();
        createCircleModalFragment.show(getFragmentManager(), "dialog");
    }

    public void setCircleAdapter (CircleAdapter circleAdapter) {
        this.adapter = circleAdapter;
    }

    public void getCircles(View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String sessionKey = "com.cherami.cherami.token";
        String token = prefs.getString(sessionKey, null);
        String userKey = "com.cherami.cherami.username";
        String username = prefs.getString(userKey, null);
        RequestParams params = new RequestParams();
        params.put("user", username);
        final View view2 = view;

        client.addHeader("Authorization", token);
        client.get(getActivity().getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "circles",
                   params, new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                System.out.println("Starting GET Circles Request");
            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;
                try {
                    responseText = new JSONObject(new String(responseBody)).getString("results");
                    JSONArray y = new JSONArray(responseText);
                    System.out.println(y.toString());
                    Circle circle_data[] = new Circle[y.length()];
                    for (int x = 0; x < y.length(); x++) {

                        circle_data[x] = new Circle(new JSONObject(y.get(x).toString()).getString("name"),
                                                    new JSONObject(y.get(x).toString()).getString("owner"),
                                                    processDate(new JSONObject(y.get(x).toString()).getString("created")),
                                                    new JSONObject(y.get(x).toString()).getString("visibility"));
                    }

                    CircleAdapter adapter = new CircleAdapter(getActivity(),
                            R.layout.circle_item_row, circle_data);
                    Circles.this.setCircleAdapter(adapter);
                    circleList = (ListView) view2.findViewById(R.id.circleList);

                    circleList.setAdapter(adapter);
                } catch (JSONException j) {
                    System.out.println(j);
                }
                filterCircles();
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

    // TODO: Rename method, update argument and hook method into UI event
    public void onButtonPressed(Uri uri) {
        if (mListener != null) {
            mListener.onFragmentInteraction(uri);
        }
    }

    public String processDate(String date){
        return date.substring(0, date.lastIndexOf("T"));
    }

    @Override
    public void onAttach(Activity activity) {
        super.onAttach(activity);
        try {
            mListener = (OnFragmentInteractionListener) activity;
        } catch (ClassCastException e) {
            throw new ClassCastException(activity.toString()
                    + " must implement OnFragmentInteractionListener");
        }
    }

    @Override
    public void onDetach() {
        super.onDetach();
        mListener = null;
    }

    public interface OnFragmentInteractionListener {
        // TODO: Update argument type and name
        public void onFragmentInteraction(Uri uri);
    }

}
