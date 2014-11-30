package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

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
        holder.txtName.setText(circle.name);
        holder.txtOwner.setText(circle.owner);
        holder.txtDate.setText(circle.date);

        return row;
    }

    static class CircleHolder {
        TextView txtName;
        TextView txtOwner;
        TextView txtDate;
    }
}
