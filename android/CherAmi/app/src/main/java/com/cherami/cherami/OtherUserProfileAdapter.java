package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import org.json.JSONException;

/**
 * Created by Geoff on 11/30/2014.
 */
public class OtherUserProfileAdapter extends ArrayAdapter<OtherUserCircle> {


    Context context;
    int layoutResourceId;
    OtherUserCircle[] data = null;

    public OtherUserProfileAdapter(Context context, int layoutResourceId, OtherUserCircle[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }


    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row = convertView;
        OtherUserCircleHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new OtherUserCircleHolder();
            holder.txtName = (TextView) row.findViewById(R.id.txtName);
            holder.txtOwner = (TextView) row.findViewById(R.id.txtOwner);
            holder.txtDate = (TextView) row.findViewById(R.id.txtDate);

            row.setTag(holder);
        } else {
            holder = (OtherUserCircleHolder) row.getTag();
        }

        OtherUserCircle otherUserCircle = data[position];
        try {
            holder.txtName.setText(otherUserCircle.circle.getString("name"));
            holder.txtOwner.setText(otherUserCircle.circle.getString("owner"));
            holder.txtDate.setText(otherUserCircle.circle.getString("created"));
        } catch (JSONException e) {
            e.printStackTrace();
        }

        return row;
    }

    @Override
    public OtherUserCircle getItem(int position){
        return data[position];
    }

    static class OtherUserCircleHolder {
        TextView txtName;
        TextView txtOwner;
        TextView txtDate;
    }
}
