package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.graphics.Color;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import org.json.JSONException;
import org.w3c.dom.Text;

public class CircleAdapter extends ArrayAdapter<Circle> {

    Context context;
    int layoutResourceId;
    Circle[] data = null;

    public CircleAdapter(Context context, int layoutResourceId, Circle[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }

    public String processDate(String date){
        return date.substring(0, date.lastIndexOf("T"));
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row = convertView;
        CircleHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new CircleHolder();
            holder.txtName = (TextView) row.findViewById(R.id.txtName);
            holder.txtOwner = (TextView) row.findViewById(R.id.txtOwner);
            holder.txtDate = (TextView) row.findViewById(R.id.txtDate);

            row.setTag(holder);
        } else {
            holder = (CircleHolder) row.getTag();
        }

        Circle circle = data[position];

        try {
            if(circle.circle.getString("visibility").equals("private")){
                row.setBackgroundColor(Color.parseColor("#4cc1f0"));
                holder.txtName.setTextColor(Color.parseColor("#ffffff"));
                holder.txtOwner.setTextColor(Color.parseColor("#ffffff"));
                holder.txtDate.setTextColor(Color.parseColor("#ffffff"));
            } else {
                row.setBackgroundColor(Color.parseColor("#f0f0f0f0"));
                holder.txtName.setTextColor(Color.parseColor("#000000"));
                holder.txtOwner.setTextColor(Color.parseColor("#000000"));
                holder.txtDate.setTextColor(Color.parseColor("#000000"));
            }
            holder.txtName.setText(circle.circle.getString("name"));
            holder.txtOwner.setText(circle.circle.getString("owner"));
            holder.txtDate.setText(processDate(circle.circle.getString("created")));
        } catch (JSONException e) {
            e.printStackTrace();
        }

        return row;
    }

    static class CircleHolder {
        TextView txtName;
        TextView txtOwner;
        TextView txtDate;
    }

    @Override
    public Circle getItem(int position){
        return data[position];
    }

    public Circle [] getData () {
        return this.data;
    }
}
