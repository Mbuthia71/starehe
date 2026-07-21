import { View, Text, ScrollView, StyleSheet, TouchableOpacity } from 'react-native';
import { useColorScheme } from 'react-native';
import { useEffect, useState } from 'react';
import { Colors } from '@/constants/theme';
import { isAuthenticated, clearToken } from '@/lib/api';
import { useRouter } from 'expo-router';

export default function HomeScreen() {
  const scheme = useColorScheme();
  const colors = Colors[scheme === 'unspecified' ? 'light' : scheme];
  const router = useRouter();
  const [authChecked, setAuthChecked] = useState(false);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    const isAuth = await isAuthenticated();
    if (!isAuth) {
      router.replace('/auth');
    }
    setAuthChecked(true);
  };

  const handleSignOut = async () => {
    await clearToken();
    router.replace('/auth');
  };

  if (!authChecked) {
    return null;
  }

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={[styles.greeting, { color: colors.text }]}>
            Welcome back, John
          </Text>
          <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
            Old Starehian Society
          </Text>
        </View>

        <View style={styles.content}>
          <Text style={[styles.sectionTitle, { color: colors.text }]}>
            Quick Actions
          </Text>
          
          <View style={styles.quickActions}>
            <TouchableOpacity
              style={[styles.actionCard, { backgroundColor: colors.card, borderColor: colors.border }]}
              onPress={() => router.push('/directory')}
            >
              <View style={[styles.actionIcon, { backgroundColor: '#4F46E5' }]} />
              <Text style={[styles.actionLabel, { color: colors.text }]}>Directory</Text>
            </TouchableOpacity>

            <TouchableOpacity
              style={[styles.actionCard, { backgroundColor: colors.card, borderColor: colors.border }]}
              onPress={() => router.push('/jobs')}
            >
              <View style={[styles.actionIcon, { backgroundColor: '#10B981' }]} />
              <Text style={[styles.actionLabel, { color: colors.text }]}>Jobs</Text>
            </TouchableOpacity>

            <TouchableOpacity
              style={[styles.actionCard, { backgroundColor: colors.card, borderColor: colors.border }]}
              onPress={() => router.push('/tenders')}
            >
              <View style={[styles.actionIcon, { backgroundColor: '#F59E0B' }]} />
              <Text style={[styles.actionLabel, { color: colors.text }]}>Tenders</Text>
            </TouchableOpacity>

            <TouchableOpacity
              style={[styles.actionCard, { backgroundColor: colors.card, borderColor: colors.border }]}
              onPress={() => router.push('/business')}
            >
              <View style={[styles.actionIcon, { backgroundColor: '#8B5CF6' }]} />
              <Text style={[styles.actionLabel, { color: colors.text }]}>Business</Text>
            </TouchableOpacity>
          </View>

          <Text style={[styles.sectionTitle, { color: colors.text }]}>
            Recent Activity
          </Text>

          {[1, 2, 3].map((item) => (
            <View
              key={item}
              style={[styles.activityCard, { backgroundColor: colors.card, borderColor: colors.border }]}
            >
              <View style={styles.activityIcon} />
              <View style={styles.activityContent}>
                <Text style={[styles.activityTitle, { color: colors.text }]}>
                  New job posted
                </Text>
                <Text style={[styles.activityTime, { color: colors.textSecondary }]}>
                  2 hours ago
                </Text>
              </View>
            </View>
          ))}

          <TouchableOpacity
            style={[styles.signOutButton, { backgroundColor: '#EF4444' }]}
            onPress={handleSignOut}
          >
            <Text style={styles.signOutText}>Sign Out</Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    padding: 20,
    paddingTop: 60,
  },
  greeting: {
    fontSize: 28,
    fontWeight: 'bold',
  },
  subtitle: {
    fontSize: 14,
    marginTop: 4,
  },
  content: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 16,
  },
  quickActions: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginBottom: 24,
    gap: 12,
  },
  actionCard: {
    width: '48%',
    padding: 16,
    borderRadius: 12,
    borderWidth: 1,
    alignItems: 'center',
  },
  actionIcon: {
    width: 48,
    height: 48,
    borderRadius: 24,
    marginBottom: 8,
  },
  actionLabel: {
    fontSize: 14,
    fontWeight: '500',
  },
  activityCard: {
    flexDirection: 'row',
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    borderWidth: 1,
  },
  activityIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#4F46E5',
    marginRight: 12,
  },
  activityContent: {
    flex: 1,
    justifyContent: 'center',
  },
  activityTitle: {
    fontSize: 14,
    fontWeight: '500',
  },
  activityTime: {
    fontSize: 12,
    marginTop: 2,
  },
  signOutButton: {
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
    marginTop: 16,
  },
  signOutText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
