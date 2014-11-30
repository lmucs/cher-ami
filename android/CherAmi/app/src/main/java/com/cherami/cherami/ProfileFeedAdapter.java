package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.CheckBox;
import android.widget.TextView;
import android.widget.Toast;

/**
 * Created by Geoff on 11/28/2014.
 */
public class ProfileFeedAdapter extends ArrayAdapter<ProfileFeedItem> {


    Context context;
    int layoutResourceId;
    ProfileFeedItem[] data = null;

    public ProfileFeedAdapter(Context context, int layoutResourceId, ProfileFeedItem[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }


    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row = convertView;
        ProfileFeedHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new ProfileFeedHolder();
            holder.txtTitle = (TextView) row.findViewById(R.id.txtTitle);
            holder.txtContent = (TextView) row.findViewById(R.id.txtContent);
            holder.txtDate = (TextView) row.findViewById(R.id.txtDate);

            row.setTag(holder);
        } else {
            holder = (ProfileFeedHolder) row.getTag();
        }

        ProfileFeedItem profileFeedItem = data[position];
        holder.txtTitle.setText(profileFeedItem.title);
        holder.txtContent.setText(profileFeedItem.message);
        holder.txtDate.setText(profileFeedItem.date);

        return row;
    }

    static class ProfileFeedHolder {
        TextView txtTitle;
        TextView txtContent;
        TextView txtDate;
    }
}
