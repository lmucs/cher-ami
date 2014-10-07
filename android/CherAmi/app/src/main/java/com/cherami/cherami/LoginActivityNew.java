//package com.cherami.cherami;
//
//import android.app.ActionBar;
//import android.app.Activity;
//import android.app.Fragment;
//import android.app.FragmentManager;
//import android.app.FragmentTransaction;
//import android.content.Intent;
//import android.os.Bundle;
//import android.support.v13.app.FragmentPagerAdapter;
//import android.support.v4.view.ViewPager;
//import android.text.TextUtils;
//import android.view.LayoutInflater;
//import android.view.Menu;
//import android.view.MenuItem;
//import android.view.View;
//import android.view.ViewGroup;
//import android.widget.EditText;
//
//import java.util.Locale;
//
///**
// * Created by goalsman on 10/7/14.
// */
//public class LoginActivityNew extends Activity {
//    /**
//     * The {@link android.support.v4.view.PagerAdapter} that will provide
//     * fragments for each of the sections. We use a
//     * {@link android.support.v13.app.FragmentPagerAdapter} derivative, which will keep every
//     * loaded fragment in memory. If this becomes too memory intensive, it
//     * may be best to switch to a
//     * {@link android.support.v13.app.FragmentStatePagerAdapter}.
//     */
//    SectionsPagerAdapter mSectionsPagerAdapter;
//
//    /**
//     * The {@link android.support.v4.view.ViewPager} that will host the section contents.
//     */
//    ViewPager mViewPager;
//    EditText mUsername;
//    EditText mEmail;
//    EditText mPassword;
//    EditText mConfirmPassword;
//
//    @Override
//    protected void onCreate(Bundle savedInstanceState) {
//        super.onCreate(savedInstanceState);
//        setContentView(R.layout.activity_main);
//    }
//
//
//    @Override
//    public boolean onCreateOptionsMenu(Menu menu) {
//        // Inflate the menu; this adds items to the action bar if it is present.
//        getMenuInflater().inflate(R.menu.main, menu);
//        return true;
//    }
//
//    public void showLogin(View view) {
//        Intent intent = new Intent(this, FeedActivity.class);
//        startActivity(intent);
//    }
//
//    public void attemptCreateAccount(View view) {
//        View focusView = null;
//        Boolean cancel = false;
//
//        mUsername = (EditText)findViewById(R.id.username);
//        mEmail = (EditText)findViewById(R.id.email);
//        mPassword = (EditText)findViewById(R.id.password);
//        mConfirmPassword = (EditText)findViewById(R.id.confirmPassword);
//        String username = mUsername.getText().toString();
//        String email = mEmail.getText().toString();
//        String password = mPassword.getText().toString();
//        String confrimPassword = mConfirmPassword.getText().toString();
//        /* First, data sanitization: No fields should be left blank, email should have @ symbol,
//        password and Confirm password should be the same (this is done in back end)
//        Also, handle/username must be unique (also back end?)
//        Then,
//        POST to db with Handle, Email, Password, and Confirm
//        */
//
//        // Check for a valid email address.
//        if (TextUtils.isEmpty(email)) {
//            mEmail.setError(getString(R.string.error_field_required));
//            focusView = mEmail;
//            cancel = true;
//        } else if (!isEmailValid(email)) {
//            mEmail.setError(getString(R.string.error_invalid_email));
//            focusView = mEmail;
//            cancel = true;
//        }
//
//        if (TextUtils.isEmpty(password)) {
//            mPassword.setError(getString(R.string.error_field_required));
//            focusView = mPassword;
//            cancel = true;
//        } else if (!isPasswordValid(password)) {
//            mPassword.setError(getString(R.string.error_invalid_password));
//            focusView = mPassword;
//            cancel = true;
//        }
//
//        if (!confrimPassword.equals(password)) {
//            mConfirmPassword.setError(getString(R.string.error_invalid_confirm));
//            focusView = mConfirmPassword;
//            cancel = true;
//        }
//
//        if (TextUtils.isEmpty(username)) {
//            mUsername.setError(getString(R.string.error_field_required));
//            focusView = mUsername;
//            cancel = true;
//        }
//
//
//        if (cancel) {
//            //Something is wrong; don't sign up
//            focusView.requestFocus();
//        } else {
//            // Sign them up
//
//        }
//
//    }
//
//    private boolean isEmailValid(String email) {
//        //TODO: Replace this with your own logic
//        return email.contains("@");
//    }
//
//    private boolean isPasswordValid(String password) {
//        //TODO: Replace this with your own logic
//        return password.length() > 4;
//    }
//
//
//
//    @Override
//    public boolean onOptionsItemSelected(MenuItem item) {
//        // Handle action bar item clicks here. The action bar will
//        // automatically handle clicks on the Home/Up button, so long
//        // as you specify a parent activity in AndroidManifest.xml.
//        int id = item.getItemId();
//        if (id == R.id.action_settings) {
//            return true;
//        }
//        return super.onOptionsItemSelected(item);
//    }
//
//    /**
//     * A {@link android.support.v13.app.FragmentPagerAdapter} that returns a fragment corresponding to
//     * one of the sections/tabs/pages.
//     */
//    public class SectionsPagerAdapter extends FragmentPagerAdapter {
//
//        public SectionsPagerAdapter(FragmentManager fm) {
//            super(fm);
//        }
//
//        @Override
//        public Fragment getItem(int position) {
//            // getItem is called to instantiate the fragment for the given page.
//            // Return a PlaceholderFragment (defined as a static inner class below).
//            return PlaceholderFragment.newInstance(position + 1);
//        }
//
//        @Override
//        public int getCount() {
//            // Show 3 total pages.
//            return 3;
//        }
//
//        @Override
//        public CharSequence getPageTitle(int position) {
//            Locale l = Locale.getDefault();
//            switch (position) {
//                case 0:
//                    return getString(R.string.title_section1).toUpperCase(l);
//                case 1:
//                    return getString(R.string.title_section2).toUpperCase(l);
//                case 2:
//                    return getString(R.string.title_section3).toUpperCase(l);
//            }
//            return null;
//        }
//    }
//
//    /**
//     * A placeholder fragment containing a simple view.
//     */
//    public static class PlaceholderFragment extends Fragment {
//        /**
//         * The fragment argument representing the section number for this
//         * fragment.
//         */
//        private static final String ARG_SECTION_NUMBER = "section_number";
//
//        /**
//         * Returns a new instance of this fragment for the given section
//         * number.
//         */
//        public static PlaceholderFragment newInstance(int sectionNumber) {
//            PlaceholderFragment fragment = new PlaceholderFragment();
//            Bundle args = new Bundle();
//            args.putInt(ARG_SECTION_NUMBER, sectionNumber);
//            fragment.setArguments(args);
//            return fragment;
//        }
//
//        public PlaceholderFragment() {
//        }
//
//        @Override
//        public View onCreateView(LayoutInflater inflater, ViewGroup container,
//                                 Bundle savedInstanceState) {
//            View rootView = inflater.inflate(R.layout.fragment_main, container, false);
//            return rootView;
//        }
//    }
//}
