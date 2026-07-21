import { View, Text, ScrollView, StyleSheet, TouchableOpacity } from 'react-native';
import { useColorScheme } from 'react-native';
import { Colors } from '@/constants/theme';
import { getBusinessListings } from '@/lib/api';

export default function DirectoryScreen() {
  const scheme = useColorScheme();
  const colors = Colors[scheme === 'unspecified' ? 'light' : scheme];

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: colors.text }]}>Directory</Text>
          <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
            Connect with fellow Old Starehians
          </Text>
        </View>

        <View style={styles.searchContainer}>
          <Text style={[styles.searchPlaceholder, { color: colors.textSecondary }]}>
            Search alumni...
          </Text>
        </View>

        <View style={styles.content}>
          <Text style={[styles.sectionTitle, { color: colors.text }]}>
            Featured Alumni
          </Text>
          
          {/* Sample alumni cards */}
          {[1, 2, 3].map((item) => (
            <TouchableOpacity
              key={item}
              style={[styles.card, { backgroundColor: colors.card, borderColor: colors.border }]}
            >
              <View style={styles.avatar} />
              <View style={styles.cardContent}>
                <Text style={[styles.name, { color: colors.text }]}>
                  John Doe {item}
                </Text>
                <Text style={[styles.role, { color: colors.textSecondary }]}>
                  Class of 2015
                </Text>
                <Text style={[styles.location, { color: colors.textSecondary }]}>
                  Nairobi, Kenya
                </Text>
              </View>
            </TouchableOpacity>
          ))}
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
  title: {
    fontSize: 32,
    fontWeight: 'bold',
  },
  subtitle: {
    fontSize: 14,
    marginTop: 4,
  },
  searchContainer: {
    marginHorizontal: 20,
    marginBottom: 20,
    padding: 16,
    borderRadius: 12,
    borderWidth: 1,
  },
  searchPlaceholder: {
    fontSize: 14,
  },
  content: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 16,
  },
  card: {
    flexDirection: 'row',
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    borderWidth: 1,
  },
  avatar: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: '#4F46E5',
    marginRight: 12,
  },
  cardContent: {
    flex: 1,
  },
  name: {
    fontSize: 16,
    fontWeight: '600',
  },
  role: {
    fontSize: 14,
    marginTop: 2,
  },
  location: {
    fontSize: 12,
    marginTop: 4,
  },
});
